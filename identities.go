package migagoapi

import (
	"context"
	"encoding/json"
	"fmt"
)

type Identity struct {
	LocalPart            string `json:"local_part,omitempty"`
	Domain               string `json:"domain,omitempty"`
	Address              string `json:"address,omitempty"`
	Name                 string `json:"name,omitempty"`
	MaySend              bool   `json:"may_send,omitempty"`
	MayReceive           bool   `json:"may_receive,omitempty"`
	MayAccessImap        bool   `json:"may_access_imap,omitempty"`
	MayAccessPop3        bool   `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve bool   `json:"may_access_managesieve,omitempty"`
	Password             string `json:"password,mayomitempty"`
	PasswordUse          string `json:"password_use,omitempty"` // leaked info on new documentation (!!!) from https://github.com/metio/migadu-client.go/pull/9
	FooterActive         bool   `json:"footer_active,omitempty"`
	FooterPlainBody      string `json:"footer_plain_body,omitempty"`
	FooterHtmlBody       string `json:"footer_html_body,omitempty"`
}

func (c *Client) GetIdentities(ctx context.Context, localPart string) ([]Identity, error) {
	var identities struct {
		Identities []Identity `json:"Identities,omitempty"`
	}

	body, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities", localPart))
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	err = json.Unmarshal(body, &identities)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	return identities.Identities, nil
}

func (c *Client) GetIdentity(ctx context.Context, localPart, id string) (*Identity, error) {
	var identity Identity

	body, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities/%s", localPart, id))
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	return &identity, nil
}

func (c *Client) CreateIdentity(ctx context.Context, localPart string, newIdentity *Identity) (*Identity, error) {
	var identity Identity

	urlSlug := fmt.Sprintf("mailboxes/%s/identities", localPart)

	jsonBody, err := json.Marshal(newIdentity)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	body, err := c.Post(ctx, urlSlug, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	return &identity, nil
}

// Convience function to create an identity with a password.
// The password can be a custom password as passed, or left as an empty string to use the mailbox password
func (c *Client) CreateIdentityWithPassword(ctx context.Context, localPart, name, id, password string) (*Identity, error) {
	newIdentity := Identity{Name: name, LocalPart: id}
	if password == "" {
		newIdentity.Password = ""
		newIdentity.PasswordUse = "mailbox"
	} else {
		newIdentity.Password = password
		newIdentity.PasswordUse = "custom"
	}

	return c.CreateIdentity(ctx, localPart, &newIdentity)
}

// Convience function to create an identity that will not be used for authentication (i.e. login)
func (c *Client) CreateIdentityNoAuth(ctx context.Context, localPart, name, id string) (*Identity, error) {
	newIdentity := Identity{Name: name, LocalPart: id, PasswordUse: "none"}

	return c.CreateIdentity(ctx, localPart, &newIdentity)
}

func (c *Client) UpdateIdentity(ctx context.Context, localPart, id string, identityParams *Identity) (*Identity, error) {
	var identity Identity

	urlSlug := fmt.Sprintf("mailboxes/%s/identities/%s", localPart, id)

	jsonBody, err := json.Marshal(identityParams)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	body, err := c.Put(ctx, urlSlug, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	return &identity, nil
}

func (c *Client) DeleteIdentity(ctx context.Context, localPart, id string) error {
	urlSlug := fmt.Sprintf("mailboxes/%s/identities/%s", localPart, id)

	_, err := c.Delete(ctx, urlSlug)

	if err != nil {
		return fmt.Errorf("DeleteIdentity: %w", err)
	}

	return nil
}
