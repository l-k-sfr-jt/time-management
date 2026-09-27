package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

// webhookTolerance rejects a webhook whose svix-timestamp is further than
// this from now, in either direction, guarding against replay of a
// captured payload.
const webhookTolerance = 5 * time.Minute

// verifyWebhookSignature implements the Svix/Standard Webhooks signing
// scheme Clerk uses: HMAC-SHA256 over "{id}.{timestamp}.{body}", keyed by
// the base64 payload of the "whsec_"-prefixed secret. We hand-roll this
// (rather than pulling in the svix-webhooks SDK, which vendors its entire
// multi-language monorepo as a single Go module) because the algorithm is
// small, stable, and fully specified.
func verifyWebhookSignature(secret string, r *http.Request, body []byte) error {
	id := r.Header.Get("svix-id")
	timestamp := r.Header.Get("svix-timestamp")
	signatureHeader := r.Header.Get("svix-signature")
	if id == "" || timestamp == "" || signatureHeader == "" {
		return fmt.Errorf("missing svix headers")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid svix-timestamp: %w", err)
	}
	if age := time.Since(time.Unix(ts, 0)); age > webhookTolerance || age < -webhookTolerance {
		return fmt.Errorf("svix-timestamp outside tolerance")
	}

	secretBytes, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	if err != nil {
		return fmt.Errorf("invalid webhook secret encoding: %w", err)
	}

	signedContent := fmt.Sprintf("%s.%s.%s", id, timestamp, body)
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signedContent))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	for _, part := range strings.Fields(signatureHeader) {
		version, sig, found := strings.Cut(part, ",")
		if !found || version != "v1" {
			continue
		}
		if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) == 1 {
			return nil
		}
	}
	return fmt.Errorf("no matching signature")
}

type clerkWebhookPayload struct {
	Type string `json:"type"`
	Data struct {
		ID                    string `json:"id"`
		PrimaryEmailAddressID string `json:"primary_email_address_id"`
		EmailAddresses        []struct {
			ID           string `json:"id"`
			EmailAddress string `json:"email_address"`
		} `json:"email_addresses"`
	} `json:"data"`
}

func (p clerkWebhookPayload) primaryEmail() string {
	for _, e := range p.Data.EmailAddresses {
		if e.ID == p.Data.PrimaryEmailAddressID {
			return e.EmailAddress
		}
	}
	return ""
}

// ClerkWebhookHandler handles user.created/user.updated events by
// upserting the local users row (architecture/data-flows.md §5). Other
// event types are accepted (200) and ignored.
func ClerkWebhookHandler(secret string, queries *sqlcgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		if err := verifyWebhookSignature(secret, r, body); err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var payload clerkWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		switch payload.Type {
		case "user.created", "user.updated":
			if _, err := queries.UpsertUser(r.Context(), sqlcgen.UpsertUserParams{
				ClerkUserID: payload.Data.ID,
				Email:       payload.primaryEmail(),
			}); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
