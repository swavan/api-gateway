package salt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"os"
)

type API interface {
	Encrypt(text string) (string, error)
	Decrypt(cipherText string) (string, error)
	GenerateKey() (public string, private string, err error)
}

type Salt struct {
	secret string
}

func (s *Salt) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(s.secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	bytes := gcm.Seal(
		nonce,
		nonce,
		[]byte(text),
		nil)
	return string(base64.RawURLEncoding.EncodeToString(bytes)), nil
}

func (s *Salt) Decrypt(cipherText string) (string, error) {
	keyByte := []byte(s.secret)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()

	bytes, err := base64.RawStdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	nonce, ByteClean := bytes[:nonceSize], bytes[nonceSize:]
	plaintextByte, err := gcm.Open(nil, nonce, ByteClean, nil)
	return string(plaintextByte), err
}

func (s *Salt) GenerateKey() (public string, private string, err error) {
	pub, priv, errs := ed25519.GenerateKey(nil)
	if errs != nil {
		return "", "", errs
	}

	privateBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", err
	}
	// private = hex.EncodeToString(privateBytes)
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	}

	var privateKeyBuffer bytes.Buffer
	if err := pem.Encode(&privateKeyBuffer, privateKeyBlock); err != nil {
		return "", "", err
	}

	private = privateKeyBuffer.String()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", "", err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	var publicKeyBuffer bytes.Buffer
	if err := pem.Encode(&publicKeyBuffer, publicKeyBlock); err != nil {
		return "", "", err
	}
	public = publicKeyBuffer.String()
	return
}

func New(secrets ...string) API {
	secret := os.Getenv("SECRET_SALT")
	for _, s := range secrets {
		secret = s
	}
	return &Salt{secret: secret}
}
