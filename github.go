package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	gitHubSecretEnvVar = "GITHUB_APP_WEBHOOK_SECRET"
)

// Webhook parse errors
var (
	ErrHMACVerificationFailed    = errors.New("HMAC verification failed")
	ErrInvalidHTTPMethod         = errors.New("invalid HTTP Method")
	ErrMissingHubSignatureHeader = errors.New("missing X-Hub-Signature Header")
	ErrReadingPayload            = errors.New("error parsing payload")
)

func verifyGithubSignature(r *http.Request, payload []byte) error {
	secret := os.Getenv(gitHubSecretEnvVar)

	signature := r.Header.Get("X-Hub-Signature-256")
	if len(signature) == 0 {
		return ErrMissingHubSignatureHeader
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
		return ErrHMACVerificationFailed
	}

	return nil
}

func receiveGithubWebhook(r *http.Request) (*string, []byte, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		return nil, nil, ErrInvalidHTTPMethod
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, nil, ErrReadingPayload
	}

	err = verifyGithubSignature(r, payload)
	if err != nil {
		return nil, nil, err
	}

	deliveryID := r.Header.Get("X-GitHub-Delivery")

	return &deliveryID, payload, nil
}
