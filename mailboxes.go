package migagoapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const subpath = "mailboxes"

type CustomTime time.Time

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}
	*t = CustomTime(date)
	return
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"2006-01-02T15:04:05.000-0700"`)), nil
}

type Mailbox struct {
	LocalPart             string     `json:"local_part,omitempty"`
	Domain                string     `json:"domain,omitempty"`
	Address               string     `json:"address,omitempty"`
	Name                  string     `json:"name,omitempty"`
	IsInternal            bool       `json:"is_internal,omitempty"`
	MaySend               bool       `json:"may_send,omitempty"`
	MayReceive            bool       `json:"may_receive,omitempty"`
	MayAccessImap         bool       `json:"may_access_imap,omitempty"`
	MayAccessPop3         bool       `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve  bool       `json:"may_access_managesieve,omitempty"`
	PasswordMethod        string     `json:"password_method,omitempty"`
	Password              string     `json:"password,omitempty"`
	PasswordRecoveryEmail string     `json:"password_recovery_email,omitempty"`
	SpamAction            string     `json:"spam_action,omitempty"`
	SpamAggressiveness    string     `json:"spam_aggressiveness"`
	SenderDenylist        []string   `json:"sender_denylist,omitempty"`
	SenderAllowlist       []string   `json:"sender_allowlist,omitempty"`
	RecipientDenylist     []string   `json:"recipient_denylist,omitempty"`
	AutorespondActive     bool       `json:"autorespond_active,omitempty"`
	AutorespondSubject    string     `json:"autorespond_subject,omitempty"`
	AutorespondBody       string     `json:"autorespond_body,omitempty"`
	AutorespondExpiresOn  string     `json:"autorespond_expires_on,omitempty"`
	FooterActive          bool       `json:"footer_active,omitempty"`
	FooterPlainBody       string     `json:"footer_plain_body,omitempty"`
	FooterHTMLBody        string     `json:"footer_html_body,omitempty"`
	Identities            []Identity `json:"identities,omitempty"`
}

// Get all mailboxes on the domain associated with the client
func (c *Client) GetMailboxes(ctx context.Context) ([]Mailbox, error) {
	var mailboxes struct {
		Mailboxes []Mailbox `json:"mailboxes,omitempty"`
	}

	body, err := c.Get(ctx, "mailboxes")
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	err = json.Unmarshal(body, &mailboxes)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	return mailboxes.Mailboxes, nil
}

// Get mailbox local_part at domain associated with the client
func (c *Client) GetMailbox(ctx context.Context, localPart string) (*Mailbox, error) {
	var mailbox Mailbox

	url_slug := fmt.Sprintf("mailboxes/%s", localPart)

	body, err := c.Get(ctx, url_slug)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	err = json.Unmarshal(body, &mailbox)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	return &mailbox, nil
}

// Create mailbox using Mailbox object
func (c *Client) CreateMailbox(ctx context.Context, newMailbox *Mailbox) (*Mailbox, error) {
	var mailbox Mailbox

	mailboxBody, err := json.Marshal(newMailbox)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: When marshaling, %w", err)
	}

	body, err := c.Post(ctx, "mailboxes", mailboxBody)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: When posting, %w", err)
	}

	err = json.Unmarshal(body, &mailbox)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: When unmarshaling, %w", err)
	}

	return &mailbox, nil
}

// Convience function to create a mailbox with a password set
func (c *Client) CreateMailboxWithPassword(
	ctx context.Context, name, localPart, password string, isInternal bool) (*Mailbox, error) {

	newMailbox := Mailbox{
		Name:       name,
		LocalPart:  localPart,
		Password:   password,
		IsInternal: isInternal,
	}

	return c.CreateMailbox(ctx, &newMailbox)
}

// Convience function to create a mailbox that sets the password via invitation link
func (c *Client) CreateMailboxWithInvite(
	ctx context.Context, name, localPart, passwordRecoveryEmail string) (*Mailbox, error) {
	newMailbox := Mailbox{
		Name:                  name,
		LocalPart:             localPart,
		PasswordMethod:        "invitation",
		PasswordRecoveryEmail: passwordRecoveryEmail,
	}

	return c.CreateMailbox(ctx, &newMailbox)
}

// Updates the mailbox local_part using the parametres in the provided Mailbox
// Returns the updated Mailbox as a pointer and any errors
func (c *Client) UpdateMailbox(ctx context.Context, localPart string, mailboxParams *Mailbox) (*Mailbox, error) {
	var updatedMailbox Mailbox

	urlSlug := fmt.Sprintf("mailboxes/%s", localPart)

	mailbox_body, err := json.Marshal(mailboxParams)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	body, err := c.Put(ctx, urlSlug, mailbox_body)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	err = json.Unmarshal(body, &updatedMailbox)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	return &updatedMailbox, nil
}

func (c *Client) DeleteMailbox(ctx context.Context, localPart string) error {
	urlSlug := fmt.Sprintf("mailboxes/%s", localPart)

	_, err := c.Delete(ctx, urlSlug)
	if err != nil {
		return fmt.Errorf("DeleteMailbox: %w", err)
	}

	return nil
}
