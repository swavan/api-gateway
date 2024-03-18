package secret

import "time"

type Secret struct {
	ID          string    `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"`
	Description string    `json:"description" db:"description"`
	Domain      string    `json:"domain" db:"domain"`
	IssueAt     time.Time `json:"issue_at" db:"issue_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	AlertTo     string    `json:"alert_to" db:"alert_to"`
	Modifier    string    `json:"modifier" db:"modifier"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func NewSecret() *Secret {
	return &Secret{}
}

func (t *Secret) SetID(id string) *Secret {
	t.ID = id
	return t
}

func (t *Secret) SetType(ty string) *Secret {
	t.Type = ty
	return t
}

func (t *Secret) SetDescription(description string) *Secret {
	t.Description = description
	return t
}

func (t *Secret) SetDomain(domain string) *Secret {
	t.Domain = domain
	return t
}

func (t *Secret) SetIssueAt(issueAt time.Time) *Secret {
	t.IssueAt = issueAt
	return t
}

func (t *Secret) SetExpiresAt(expiresAt time.Time) *Secret {
	t.ExpiresAt = expiresAt
	return t
}

func (t *Secret) SetAlertTo(alertTo string) *Secret {
	t.AlertTo = alertTo
	return t
}

func (t *Secret) SetModifier(modifier string) *Secret {
	t.Modifier = modifier
	return t
}

func (t *Secret) SetCreatedAt(createdAt time.Time) *Secret {
	t.CreatedAt = createdAt
	return t
}

func (t *Secret) GetID() string {
	return t.ID
}
