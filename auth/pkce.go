package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc7636

func encodeB64UrlSafe(b string) string {
	data := base64.RawURLEncoding.EncodeToString([]byte(b))

	out := data
	out = strings.Replace(out, "+", "-", -1)
	out = strings.Replace(out, "/", "_", -1)
	out = strings.Replace(out, "=", "", -1)
	return out
}

func getChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return encodeB64UrlSafe(string(hash[:]))
}

func getVerifier() (string, error) {
	randBuf := make([]byte, 32)
	_, err := rand.Read(randBuf)
	if err != nil {
		return "", err
	}
	return encodeB64UrlSafe(hex.EncodeToString(randBuf)), nil
}
