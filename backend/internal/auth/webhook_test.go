package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func sign(t *testing.T, secretB64, id, timestamp, body string) string {
	t.Helper()
	secretBytes, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		t.Fatalf("decode secret: %v", err)
	}
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(fmt.Sprintf("%s.%s.%s", id, timestamp, body)))
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func newSignedRequest(t *testing.T, secret, secretB64, body string, ts time.Time) *http.Request {
	t.Helper()
	id := "msg_test123"
	timestamp := strconv.FormatInt(ts.Unix(), 10)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/clerk", nil)
	req.Header.Set("svix-id", id)
	req.Header.Set("svix-timestamp", timestamp)
	req.Header.Set("svix-signature", sign(t, secretB64, id, timestamp, body))
	return req
}

func TestVerifyWebhookSignature(t *testing.T) {
	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		t.Fatalf("generate secret: %v", err)
	}
	secretB64 := base64.StdEncoding.EncodeToString(rawSecret)
	secret := "whsec_" + secretB64
	body := `{"type":"user.created","data":{"id":"user_abc"}}`

	t.Run("valid signature is accepted", func(t *testing.T) {
		req := newSignedRequest(t, secret, secretB64, body, time.Now())
		if err := verifyWebhookSignature(secret, req, []byte(body)); err != nil {
			t.Fatalf("expected valid signature, got error: %v", err)
		}
	})

	t.Run("tampered body is rejected", func(t *testing.T) {
		req := newSignedRequest(t, secret, secretB64, body, time.Now())
		tampered := body + "extra"
		if err := verifyWebhookSignature(secret, req, []byte(tampered)); err == nil {
			t.Fatal("expected tampered body to be rejected")
		}
	})

	t.Run("wrong secret is rejected", func(t *testing.T) {
		req := newSignedRequest(t, secret, secretB64, body, time.Now())
		otherSecretBytes := make([]byte, 32)
		_, _ = rand.Read(otherSecretBytes)
		wrongSecret := "whsec_" + base64.StdEncoding.EncodeToString(otherSecretBytes)
		if err := verifyWebhookSignature(wrongSecret, req, []byte(body)); err == nil {
			t.Fatal("expected wrong secret to be rejected")
		}
	})

	t.Run("stale timestamp is rejected", func(t *testing.T) {
		req := newSignedRequest(t, secret, secretB64, body, time.Now().Add(-1*time.Hour))
		if err := verifyWebhookSignature(secret, req, []byte(body)); err == nil {
			t.Fatal("expected stale timestamp to be rejected")
		}
	})

	t.Run("missing headers are rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/clerk", nil)
		if err := verifyWebhookSignature(secret, req, []byte(body)); err == nil {
			t.Fatal("expected missing svix headers to be rejected")
		}
	})
}
