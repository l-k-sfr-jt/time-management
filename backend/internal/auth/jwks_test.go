package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startJWKSServer(t *testing.T, kid string, pub *rsa.PublicKey) *httptest.Server {
	t.Helper()
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big3Bytes(pub.E))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jwksResponse{
			Keys: []jwk{{Kty: "RSA", Kid: kid, N: n, E: e}},
		})
	}))
	t.Cleanup(server.Close)
	return server
}

func big3Bytes(e int) []byte {
	// Standard 65537 exponent as 3 bytes; sufficient for these tests since
	// rsa.GenerateKey always uses 65537.
	return []byte{byte(e >> 16), byte(e >> 8), byte(e)}
}

func TestJWKSCacheFetchesAndCachesKeys(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	server := startJWKSServer(t, "test-kid", &priv.PublicKey)

	cache := NewJWKSCache(server.URL, time.Minute)
	key, err := cache.Key(t.Context(), "test-kid")
	if err != nil {
		t.Fatalf("expected key, got error: %v", err)
	}
	if key.N.Cmp(priv.PublicKey.N) != 0 {
		t.Fatal("returned key does not match expected public key")
	}

	if _, err := cache.Key(t.Context(), "unknown-kid"); err == nil {
		t.Fatal("expected error for unknown kid")
	}
}
