package salt

import (
	"bytes"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

type Hash struct {
	Value  []byte
	Secret []byte
}

func NewHash() *Hash {
	return &Hash{}
}

func (h *Hash) SetValue(value []byte) *Hash {
	h.Value = value
	return h
}

func (h *Hash) SetSecret(secret []byte) *Hash {
	h.Secret = secret
	return h
}

type PasswordManagerConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func (a *PasswordManagerConfig) SetTime(time uint32) *PasswordManagerConfig {
	a.time = time
	return a
}

func (a *PasswordManagerConfig) SetMemory(memory uint32) *PasswordManagerConfig {
	a.memory = memory
	return a
}

func (a *PasswordManagerConfig) SetThreads(threads uint8) *PasswordManagerConfig {
	a.threads = threads
	return a
}

func (a *PasswordManagerConfig) SetKeyLen(keyLen uint32) *PasswordManagerConfig {
	a.keyLen = keyLen
	return a
}

func (a *PasswordManagerConfig) SetSaltLen(saltLen uint32) *PasswordManagerConfig {
	a.saltLen = saltLen
	return a
}

func NewPasswordManagerConfig() *PasswordManagerConfig {
	return &PasswordManagerConfig{}
}

type PasswordManager struct {
	config *PasswordManagerConfig
	secret []byte
}

func NewPasswordManager(config *PasswordManagerConfig) (pm *PasswordManager, err error) {
	pm = &PasswordManager{config: config}
	pm.secret, err = GenerateRandomSecret(config.saltLen)
	return pm, err
}

func (a *PasswordManager) SetSecret(secret []byte) *PasswordManager {
	a.secret = secret
	return a
}

func GenerateRandomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (a *PasswordManager) GenerateHash(password []byte) (*Hash, error) {
	h := NewHash()
	var err error
	if len(a.secret) == 0 {
		a.secret, err = GenerateRandomSecret(a.config.saltLen)
	}
	if err != nil {
		return h, err
	}
	return h.SetValue(argon2.IDKey(
		password,
		a.secret,
		a.config.time,
		a.config.memory,
		a.config.threads,
		a.config.keyLen)), nil
}

func (a *PasswordManager) Compare(hash, password []byte) error {
	salt, err := a.GenerateHash(password)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, salt.Value) {
		return errors.New("hash doesn't match")
	}
	return nil
}
