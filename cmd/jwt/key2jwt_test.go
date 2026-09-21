package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
)

func testPEMKey(t *testing.T) []byte {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

func assertJWT(t *testing.T, jwt string) {
	t.Helper()
	if parts := strings.Split(jwt, "."); len(parts) != 3 {
		t.Fatalf("expected a JWT with 3 parts, got %q", jwt)
	}
}

func TestDetectKeyFormat(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    keyFormat
		wantErr error
	}{
		{name: "pem", key: "-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----\n", want: keyFormatPEM},
		{name: "pem with leading whitespace", key: "\n\n  -----BEGIN PRIVATE KEY-----\n", want: keyFormatPEM},
		{name: "json", key: `{"type":"application"}`, want: keyFormatJSON},
		{name: "json with leading whitespace", key: "\n\t {\"type\":\"serviceaccount\"}", want: keyFormatJSON},
		{name: "empty", key: "", want: keyFormatUnknown, wantErr: errUnknownKeyFormat},
		{name: "garbage", key: "not a key", want: keyFormatUnknown, wantErr: errUnknownKeyFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectKeyFormat([]byte(tt.key))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("format = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateJWT_PEM(t *testing.T) {
	audience = "https://example.com"
	key := testPEMKey(t)

	jwt, err := generateJWT(key, "client_id")
	if err != nil {
		t.Fatalf("generateJWT: %v", err)
	}
	assertJWT(t, jwt)

	if _, err := generateJWT(key, ""); err == nil {
		t.Fatal("expected an error when issuer is missing for a PEM key")
	}
}

func TestGenerateJWT_ApplicationJSON(t *testing.T) {
	audience = "https://example.com"
	key, err := json.Marshal(map[string]string{
		"type":     "application",
		"keyId":    "key-id",
		"key":      string(testPEMKey(t)),
		"clientId": "client-id",
	})
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}

	jwt, err := generateJWT(key, "")
	if err != nil {
		t.Fatalf("generateJWT: %v", err)
	}
	assertJWT(t, jwt)
}

func TestGenerateJWT_Unknown(t *testing.T) {
	if _, err := generateJWT([]byte("garbage"), "issuer"); !errors.Is(err, errUnknownKeyFormat) {
		t.Fatalf("error = %v, want %v", err, errUnknownKeyFormat)
	}
}
