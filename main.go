package securebytes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// SecureBytes keeps encryption key and serializer.
type SecureBytes struct {
	aesKey         []byte
	additionalData []byte
	Serializer     Serializer
}

// New returns a new SecureBytes with JSONSerializer.
// `key` should provide 256 bits entropy, so if you are using random
// alphanumeric characters it should have a length of at least 50 characters.
func New(key []byte) *SecureBytes {
	hash := sha256.Sum256(key)
	return &SecureBytes{
		aesKey:         hash[:24],
		additionalData: hash[24:],
		Serializer:     JSONSerializer{},
	}
}

// Encrypt encrypts data using AES-192 with authenticated encryption,
// to avoid additional signing.
//	https://en.wikipedia.org/wiki/Authenticated_encryption
func (sb *SecureBytes) Encrypt(input interface{}) ([]byte, error) {
	var data bytes.Buffer
	err := sb.Serializer.Encode(&data, input)
	if err != nil {
		return nil, err
	}
	return sb.RawEncrypt(data.Bytes())
}

// RawEncrypt can encrypt only bytes, doesn't do serialization
func (sb *SecureBytes) RawEncrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(sb.aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, sb.additionalData)
	return ciphertext, nil
}

// Decrypt decrypts data encrypted by Encrypt.
func (sb *SecureBytes) Decrypt(data []byte, output interface{}) error {
	plaintext, err := sb.RawDecrypt(data)
	if err != nil {
		return err
	}
	return sb.Serializer.Decode(bytes.NewBuffer(plaintext), output)
}

// RawDecrypt decrypts data encrypted by RawEncrypt.
func (sb *SecureBytes) RawDecrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(sb.aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("data smaller than nonce")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, sb.additionalData)
}

// EncryptToBase64 encrypts input and converts to base64 string
func (sb *SecureBytes) EncryptToBase64(input interface{}) (string, error) {
	ciphertext, err := sb.Encrypt(input)
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

// RawEncryptToBase64 encrypts bytes and converts to base64 string
func (sb *SecureBytes) RawEncryptToBase64(data []byte) (string, error) {
	ciphertext, err := sb.RawEncrypt(data)
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

// DecryptBase64 decrypts base64 string encrypted with EncryptToBase64
func (sb *SecureBytes) DecryptBase64(b64 string, output interface{}) error {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return sb.Decrypt(data, output)
}

// RawDecryptBase64 decrypts base64 string encrypted with RawEncryptToBase64
func (sb *SecureBytes) RawDecryptBase64(b64 string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return sb.RawDecrypt(data)
}
