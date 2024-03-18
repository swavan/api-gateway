package key

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/swavan.io/gateway/pkg/authentication/salt"
)

type KeyCache map[string]*Key

type KeyManagerAPI interface {
	Salt() salt.API
	Key() KeyAPI
	CreateKey(ctx context.Context, For string) (*Key, error)
	FetchKey(ctx context.Context, useFor string) (*Key, error)
}

type KeyManager struct {
	key      KeyAPI
	salt     salt.API
	cache    KeyCache
	resource sync.RWMutex
}

// Key implements KeyManagerAPI.
func (km *KeyManager) Key() KeyAPI {
	return km.key
}

// Salt implements KeyManagerAPI.
func (km *KeyManager) Salt() salt.API {
	return km.salt
}

func (km *KeyManager) Cache() KeyCache {
	km.resource.RLock()
	defer km.resource.RUnlock()
	return km.cache
}

func (km *KeyManager) SetCache(userFor string, cache *Key) {
	km.resource.Lock()
	defer km.resource.Unlock()
	km.cache[userFor] = cache
}

func New(database *sqlx.DB, cfg *Config, sec string) (KeyManagerAPI, error) {
	km := &KeyManager{
		key:  NewKeyService(database, cfg),
		salt: salt.New(sec),
	}

	if err := km.key.Migration(context.Background()); err != nil {
		return nil, err
	}

	for _, useFor := range cfg.GenerateKey.UsrFor {
		key, err := km.FetchKey(context.Background(), useFor)
		if err != nil {
			return nil, err
		}
		km.SetCache(useFor, key)
	}
	return km, nil
}

func (km *KeyManager) CreateKey(ctx context.Context, For string) (*Key, error) {
	publicKey, privateKey, err := km.salt.GenerateKey()
	if err != nil {
		return nil, err
	}
	encryptedPrivateKey, err := km.salt.Encrypt(privateKey)
	if err != nil {
		return nil, err
	}

	k := NewKey().
		SetUseFor(For).
		SetPrivateKey(encryptedPrivateKey).
		SetPublicKey(publicKey)

	if err := km.key.Save(ctx, k, true); err != nil {
		return nil, err
	}

	return km.FetchKey(ctx, For)
}

func (km *KeyManager) FetchKey(ctx context.Context, useFor string) (*Key, error) {
	km.resource.RLock()
	defer km.resource.RUnlock()
	if key, ok := km.cache[useFor]; ok {
		return key, nil
	}

	rawKey, err := km.key.Fetch(ctx, useFor)
	if err != nil {
		// No key found so generate a new one.
		if err == sql.ErrNoRows {
			return km.CreateKey(ctx, useFor)
		}
		return nil, err
	}
	decryptedPrivateKey, err := km.salt.Decrypt(rawKey.PrivateKey)
	if err != nil {
		return nil, err
	}
	return rawKey.SetPrivateKey(decryptedPrivateKey), nil
}

type Key struct {
	ID         int64  `json:"id" db:"id"`
	UseFor     string `json:"use_for" db:"use_for"`
	PrivateKey string `json:"private_key" db:"private"`
	PublicKey  string `json:"public_key" db:"public"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

func NewKey() *Key {
	return &Key{}
}

func (k *Key) SetID(id int64) *Key {
	k.ID = id
	return k
}

func (k *Key) SetUseFor(useFor string) *Key {
	k.UseFor = useFor
	return k
}

func (k *Key) SetPrivateKey(privateKey string) *Key {
	k.PrivateKey = privateKey
	return k
}

func (k *Key) SetPublicKey(publicKey string) *Key {
	k.PublicKey = publicKey
	return k
}

func (k *Key) IsEmpty() bool {
	return k.ID == 0
}

type KeyAPI interface {
	Migration(context.Context) error
	All(context.Context) ([]Key, error)
	Find(ctx context.Context, id int64) (*Key, error)
	Fetch(ctx context.Context, useFor string) (*Key, error)
	Save(ctx context.Context, key *Key, force bool) error
	Expired(ctx context.Context, id int64) error
}

type KeyService struct {
	dep *sqlx.DB
	cfg *Config
}

// Fetch implements KeyAPI.
func (ks *KeyService) Fetch(ctx context.Context, useFor string) (*Key, error) {
	key := new(Key)
	err := ks.dep.
		GetContext(
			ctx,
			key,
			ks.cfg.Scripts.GetByUseFor,
			useFor)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return key, err
}

// All implements KeyAPI.
func (ks *KeyService) All(context.Context) ([]Key, error) {
	var keys []Key
	err := ks.dep.Select(&keys, ks.cfg.Scripts.FetchAll)
	return keys, err
}

// Expired implements KeyAPI.
func (ks *KeyService) Expired(ctx context.Context, id int64) error {
	_, err := ks.dep.
		ExecContext(
			ctx,
			ks.cfg.Scripts.DeleteByID,
			id)
	return err
}

// Find implements KeyAPI.
func (ks *KeyService) Find(ctx context.Context, id int64) (*Key, error) {
	key := new(Key)
	err := ks.dep.
		GetContext(ctx, key, ks.cfg.Scripts.GetByID, id)
	return key, err
}

// Migration implements KeyAPI.
func (ks *KeyService) Migration(ctx context.Context) error {
	if !ks.cfg.Migration.Run {
		return nil
	}
	for _, script := range ks.cfg.Migration.Scripts {
		_, err := ks.dep.ExecContext(ctx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

// Save implements KeyAPI.
func (ks *KeyService) Save(ctx context.Context, key *Key, force bool) error {
	if force {
		_, err := ks.dep.ExecContext(
			ctx,
			ks.cfg.Scripts.DeleteByUseFor, key.UseFor)
		if err != nil {
			return err
		}
	}
	_, err := ks.dep.ExecContext(
		ctx,
		ks.cfg.Scripts.Save,
		key.UseFor,
		key.PrivateKey,
		key.PublicKey,
	)
	return err
}

func NewKeyService(dep *sqlx.DB, cfg *Config) KeyAPI {
	return &KeyService{
		dep: dep,
		cfg: cfg.SetDefaultIfEmpty()}
}
