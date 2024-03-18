package authentication

import (
	"crypto"
	"crypto/ed25519"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/swavan.io/gateway/pkg/authentication/domain"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type IDTokenClaims struct {
	Subject       string `json:"sub,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	PreferredName string `json:"preferred_username,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	Name          string `json:"name,omitempty"`
	Username      string `json:"username,omitempty"`
	IIS           string `json:"iis,omitempty"`
	Type          string `json:"typ,omitempty"`
	Application   string `json:"azp,omitempty"`
	SessionState  string `json:"session_state,omitempty"`
	AuthTime      int64  `json:"auth_time,omitempty"`
	Audience      string `json:"aud,omitempty"`
}

type TokenHeader struct {
	ID         string
	Subject    string
	Audience   string
	Issuer     string
	Expiration time.Time
	NotBefore  time.Time
	IssuedAt   time.Time
}

func NewTokenHeader() *TokenHeader {
	return &TokenHeader{
		ID:         uuid.New().String(),
		NotBefore:  time.Now(),
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(time.Minute * 10),
		Issuer:     "haas",
	}
}

func (t *TokenHeader) SetID(id string) *TokenHeader {
	t.ID = id
	return t
}

func (t *TokenHeader) SetSubject(subject string) *TokenHeader {
	t.Subject = subject
	return t
}

func (t *TokenHeader) SetAudience(audience string) *TokenHeader {
	t.Audience = audience
	return t
}

func (t *TokenHeader) SetIssuer(issuer string) *TokenHeader {
	t.Issuer = issuer
	return t
}

func (t *TokenHeader) SetExpiresAt(data time.Time) *TokenHeader {
	t.Expiration = data
	return t
}

func (t *TokenHeader) SetIssuedAt(data time.Time) *TokenHeader {
	t.IssuedAt = data
	return t
}

func (t *TokenHeader) SetNotBefore(data time.Time) *TokenHeader {
	t.NotBefore = data
	return t
}

type Claims struct {
	Subject           string        `json:"sub,omitempty"`
	Username          string        `json:"username,omitempty"`
	PreferredUsername string        `json:"preferred_username,omitempty"`
	Name              string        `json:"name,omitempty"`
	GivenName         string        `json:"given_name,omitempty"`
	FamilyName        string        `json:"family_name,omitempty"`
	Email             string        `json:"email,omitempty"`
	EmailVerified     bool          `json:"email_verified,omitempty"`
	Domain            domain.Domain `json:"domain,omitempty"`
	Roles             []string      `json:"roles,omitempty"`
}

func FromIDClaims(claims IDTokenClaims) *Claims {
	cl := NewClaims().
		SetEmail(claims.Email).
		SetEmailVerified(claims.EmailVerified).
		SetName(claims.Name).
		SetFamilyName(claims.FamilyName).
		SetGivenName(claims.GivenName).
		SetPreferredUsername(claims.PreferredName).
		SetUsername(
			claims.Username,
			claims.PreferredName,
			claims.Email)

	return cl
}

func FromPasetoJSON(claims paseto.JSONToken) *Claims {
	clm := NewClaims().
		SetEmail(claims.Get("email")).
		SetEmailVerified(claims.Get("email_verified") == "true").
		SetName(claims.Get("name")).
		SetDomain(&domain.Domain{
			ID:   claims.Get("did"),
			Name: claims.Get("domain"),
		}).
		SetFamilyName(claims.Get("family_name")).
		SetGivenName(claims.Get("given_name")).
		SetPreferredUsername(claims.Get("preferred_username")).
		SetUsername(
			claims.Get("username"),
			claims.Get("preferred_username"),
			claims.Get("email"))
	return clm
}

func NewClaims() *Claims {
	return &Claims{}
}

func (t *Claims) SetUsername(usernames ...string) *Claims {
	for _, u := range usernames {
		if len(strings.TrimSpace(u)) > 0 {
			t.Username = u
		}
	}
	return t
}

func (t *Claims) SetPreferredUsername(preferredUsername string) *Claims {
	t.PreferredUsername = preferredUsername
	return t
}

func (t *Claims) SetSubject(subject string) *Claims {
	t.Subject = subject
	return t
}

func (t *Claims) SetName(name string) *Claims {
	t.Name = name
	return t
}

func (t *Claims) IsSuperUser() bool {
	return slices.Contains(t.Roles, os.Getenv("SUPER_USER_ROLE"))
}

func (t *Claims) SetGivenName(givenName string) *Claims {
	t.GivenName = givenName
	return t
}

func (t *Claims) SetFamilyName(familyName string) *Claims {
	t.FamilyName = familyName
	return t
}

func (t *Claims) SetEmail(email string) *Claims {
	t.Email = email
	return t
}

func (t *Claims) SetEmailVerified(emailVerified bool) *Claims {
	t.EmailVerified = emailVerified
	return t
}

func (t *Claims) SetDomain(domain *domain.Domain) *Claims {
	if domain != nil {
		t.Domain = *domain
	}
	return t
}

func (t *Claims) build(tkn *TokenHeader) paseto.JSONToken {
	token := paseto.JSONToken{
		Jti:        tkn.ID,
		Subject:    tkn.Subject,
		Audience:   (t.Domain.Name),
		Issuer:     tkn.Issuer,
		Expiration: tkn.Expiration,
		NotBefore:  tkn.NotBefore,
		IssuedAt:   tkn.IssuedAt,
	}

	token.Set("username", t.Username)
	token.Set("roles", strings.Join(t.Roles, ","))
	token.Set("preferred_username", t.PreferredUsername)
	token.Set("given_name", t.GivenName)
	token.Set("family_name", t.FamilyName)
	token.Set("email", t.Email)
	token.Set("did", t.Domain.ID)
	token.Set("domain", t.Domain.Name)
	token.Set("email_verified", fmt.Sprint(t.EmailVerified))
	return token
}

func (t Claims) GenerateSymmetric(secret string, sourceInfo string, header *TokenHeader) (string, error) {
	clm := t.build(header)
	return paseto.NewV2().
		Encrypt(
			[]byte(secret),
			clm,
			sourceInfo,
		)
}

func ParseED25519PrivateKey(private string) (crypto.PrivateKey, error) {
	type ed25519PrivKey struct {
		Version          int
		ObjectIdentifier struct {
			ObjectIdentifier asn1.ObjectIdentifier
		}
		PrivateKey []byte
	}
	privateBlock, _ := pem.Decode([]byte(private))
	var asn1PrivKey ed25519PrivKey
	_, err := asn1.Unmarshal(privateBlock.Bytes, &asn1PrivKey)
	if err != nil {
		return []byte{}, err
	}
	return ed25519.NewKeyFromSeed(asn1PrivKey.PrivateKey[2:]), nil
}

func ParseED25519PublicKey(publicKeyPem string) (crypto.PublicKey, error) {
	type ed25519PublicKey struct {
		OBjectIdentifier struct {
			ObjectIdentifier asn1.ObjectIdentifier
		}
		PublicKey asn1.BitString
	}
	var block *pem.Block
	block, _ = pem.Decode([]byte(publicKeyPem))

	var asn1PubKey ed25519PublicKey
	asn1.Unmarshal(block.Bytes, &asn1PubKey)

	return ed25519.PublicKey(asn1PubKey.PublicKey.Bytes), nil
}

func (t Claims) GenerateAsymmetric(private string, sourceInfo string, header *TokenHeader) (string, error) {
	privateBytes, err := ParseED25519PrivateKey(private)
	if err != nil {
		return "", err
	}
	token := t.build(header)
	return paseto.NewV2().
		Sign(
			privateBytes,
			token,
			sourceInfo,
		)
}

func ParseAsymmetricToken(token string, key string) (*Claims, string, error) {
	publicKey, err := ParseED25519PublicKey(key)
	if err != nil {
		return &Claims{}, "", err
	}
	var claims paseto.JSONToken
	var footer string
	err = paseto.NewV2().Verify(
		token,
		publicKey,
		&claims,
		&footer)
	return FromPasetoJSON(claims), footer, err
}

func ParseSymmetricToken(token string, secret string) (*Claims, string, error) {
	var pastoClaims paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(
		token,
		[]byte(secret),
		&pastoClaims,
		&footer)
	return FromPasetoJSON(pastoClaims), footer, err
}
