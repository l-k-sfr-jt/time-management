package httpapi_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/l-k-sfr-jt/time-management/backend/internal/auth"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
	"github.com/l-k-sfr-jt/time-management/backend/internal/httpapi"
)

// These tests exercise the full stack (HTTP -> auth middleware -> Postgres)
// against a real database. They need DATABASE_URL pointing at a Postgres
// instance with migrations/0001_init.sql already applied, and are skipped
// otherwise so `go test ./...` still passes in an environment without one.
func requireDatabaseURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}
	return url
}

func startJWKS(t *testing.T, kid string, pub *rsa.PublicKey) string {
	t.Helper()
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{byte(pub.E >> 16), byte(pub.E >> 8), byte(pub.E)})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(
			`{"keys":[{"kty":"RSA","kid":%q,"n":%q,"e":%q}]}`, kid, n, e,
		)))
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func signSessionToken(t *testing.T, priv *rsa.PrivateKey, kid, subject, email string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub":   subject,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

type testServer struct {
	url    string
	client *http.Client
	token  string
}

func (ts *testServer) do(t *testing.T, method, path string, body any) *http.Response {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, ts.url+path, reader)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if ts.token != "" {
		req.Header.Set("Authorization", "Bearer "+ts.token)
	}
	resp, err := ts.client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func setupTestServer(t *testing.T) *testServer {
	t.Helper()
	databaseURL := requireDatabaseURL(t)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to database: %v", err)
	}
	t.Cleanup(pool.Close)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	kid := "test-kid"
	jwksURL := startJWKS(t, kid, &priv.PublicKey)

	queries := sqlcgen.New(pool)
	verifier := auth.NewVerifier(auth.NewJWKSCache(jwksURL, time.Minute), queries)

	webhookSecretB64 := base64.StdEncoding.EncodeToString([]byte("integration-test-secret-32bytes"))
	webhookSecret := "whsec_" + webhookSecretB64

	router := httpapi.NewRouter(pool, queries, verifier, webhookSecret)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	subject := "user_" + fmt.Sprintf("%d", time.Now().UnixNano())
	token := signSessionToken(t, priv, kid, subject, subject+"@example.com")

	return &testServer{url: server.URL, client: server.Client(), token: token}
}

func TestHealthz(t *testing.T) {
	ts := setupTestServer(t)
	resp := ts.do(t, http.MethodGet, "/healthz", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestUnauthenticatedRequestIsRejected(t *testing.T) {
	ts := setupTestServer(t)
	ts.token = ""
	resp := ts.do(t, http.MethodGet, "/api/v1/groups", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestGroupsAndActivityTypesLifecycle(t *testing.T) {
	ts := setupTestServer(t)

	// Create a group.
	resp := ts.do(t, http.MethodPost, "/api/v1/groups", httpapi.CreateGroupRequest{Name: "Sport"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create group: expected 201, got %d", resp.StatusCode)
	}
	var group httpapi.Group
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		t.Fatalf("decode group: %v", err)
	}
	resp.Body.Close()
	if group.Name != "Sport" {
		t.Fatalf("expected group name Sport, got %q", group.Name)
	}

	// It shows up in the list.
	resp = ts.do(t, http.MethodGet, "/api/v1/groups", nil)
	var groups []httpapi.Group
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		t.Fatalf("decode groups: %v", err)
	}
	resp.Body.Close()
	found := false
	for _, g := range groups {
		if g.ID == group.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("created group not found in list")
	}

	// Create an Activity Type under it.
	resp = ts.do(t, http.MethodPost, "/api/v1/activity-types", httpapi.CreateActivityTypeRequest{
		Name:    "Running",
		GroupID: group.ID,
	})
	if resp.StatusCode != http.StatusCreated {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		t.Fatalf("create activity type: expected 201, got %d: %s", resp.StatusCode, body.String())
	}
	var at httpapi.ActivityType
	if err := json.NewDecoder(resp.Body).Decode(&at); err != nil {
		t.Fatalf("decode activity type: %v", err)
	}
	resp.Body.Close()
	if at.GroupID != group.ID {
		t.Fatalf("expected activity type groupId %q, got %q", group.ID, at.GroupID)
	}

	// Deleting the group is blocked while the Activity Type references it (FR-1.3).
	resp = ts.do(t, http.MethodDelete, "/api/v1/groups/"+group.ID, nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 deleting a group with activity types, got %d", resp.StatusCode)
	}

	// Archiving the Activity Type is a soft delete.
	resp = ts.do(t, http.MethodPost, "/api/v1/activity-types/"+at.ID+"/archive", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("archive activity type: expected 200, got %d", resp.StatusCode)
	}
	var archived httpapi.ActivityType
	if err := json.NewDecoder(resp.Body).Decode(&archived); err != nil {
		t.Fatalf("decode archived activity type: %v", err)
	}
	resp.Body.Close()
	if archived.ArchivedAt == nil {
		t.Fatal("expected archivedAt to be set")
	}

	// Archived types are excluded from the default listing...
	resp = ts.do(t, http.MethodGet, "/api/v1/activity-types?groupId="+group.ID, nil)
	var active []httpapi.ActivityType
	json.NewDecoder(resp.Body).Decode(&active)
	resp.Body.Close()
	for _, a := range active {
		if a.ID == at.ID {
			t.Fatal("archived activity type should not appear in default listing")
		}
	}

	// ...but do appear with includeArchived=true.
	resp = ts.do(t, http.MethodGet, "/api/v1/activity-types?groupId="+group.ID+"&includeArchived=true", nil)
	var all []httpapi.ActivityType
	json.NewDecoder(resp.Body).Decode(&all)
	resp.Body.Close()
	foundArchived := false
	for _, a := range all {
		if a.ID == at.ID {
			foundArchived = true
		}
	}
	if !foundArchived {
		t.Fatal("archived activity type should appear when includeArchived=true")
	}

	// A brand new group with no activity types can be deleted.
	resp = ts.do(t, http.MethodPost, "/api/v1/groups", httpapi.CreateGroupRequest{Name: "Temp"})
	var tempGroup httpapi.Group
	json.NewDecoder(resp.Body).Decode(&tempGroup)
	resp.Body.Close()

	resp = ts.do(t, http.MethodDelete, "/api/v1/groups/"+tempGroup.ID, nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 deleting an empty group, got %d", resp.StatusCode)
	}
}

func TestClerkWebhookUpsertsUser(t *testing.T) {
	ts := setupTestServer(t)

	secretB64 := base64.StdEncoding.EncodeToString([]byte("integration-test-secret-32bytes"))

	clerkUserID := fmt.Sprintf("user_webhook_%d", time.Now().UnixNano())
	payload := fmt.Sprintf(
		`{"type":"user.created","data":{"id":%q,"primary_email_address_id":"idn_1","email_addresses":[{"id":"idn_1","email_address":"webhook@example.com"}]}}`,
		clerkUserID,
	)

	id := "msg_test"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	secretBytes, _ := base64.StdEncoding.DecodeString(secretB64)
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(fmt.Sprintf("%s.%s.%s", id, timestamp, payload)))
	signature := "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req, _ := http.NewRequest(http.MethodPost, ts.url+"/webhooks/clerk", bytes.NewReader([]byte(payload)))
	req.Header.Set("svix-id", id)
	req.Header.Set("svix-timestamp", timestamp)
	req.Header.Set("svix-signature", signature)
	resp, err := ts.client.Do(req)
	if err != nil {
		t.Fatalf("do webhook request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
