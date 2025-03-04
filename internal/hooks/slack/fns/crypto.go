// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package fns

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"io"
	"log/slog"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/hooks/slack/errors"
)

// encrypt encrypts the given plaintext using the provided key and AES-256 in Cipher Feedback mode.
// It returns the encrypted ciphertext, or an error if encryption fails.
func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return ciphertext, nil
}

// decrypt decrypts the given ciphertext using the provided key and AES-256 in Cipher Feedback mode.
// It returns the decrypted plaintext, or an error if decryption fails.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.ErrCipherText
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// generate creates a 32-byte key for AES-256 encryption.
// It uses a SHA-512 hash of the provided workspaceID,
// and then takes the first 32 bytes of the hash as the key.
func Generate(workspaceID string) []byte {
	h := sha512.New()
	h.Write([]byte(auth.Secret() + workspaceID)) // TODO - verify

	return h.Sum(nil)[:32]
}

// Reveal decodes a base64-encoded encrypted token and decrypts it using a generated key.
func Reveal(botToken, workspaceID string) (string, error) {
	// Decode the base64-encoded encrypted token.
	decoded, err := base64.StdEncoding.DecodeString(botToken)
	if err != nil {
		slog.Error("Failed to decode the token", slog.Any("e", err))
		return "", err
	}

	// Generate the same key used for encryption.
	key := Generate(workspaceID)

	// Decrypt the token.
	decrypted, err := Decrypt(decoded, key)
	if err != nil {
		slog.Error("Failed to decrypt the token", slog.Any("e", err))
		return "", err
	}

	return string(decrypted), nil
}
