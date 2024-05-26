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

func (c *Client) GetIdentities(ctx context.Context, local_part string) (*[]Identity, error) {
	var identities struct {
		Identities []Identity `json:"Identities,omitempty"`
	}

	body, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities", local_part))
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	err = json.Unmarshal(body, &identities)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	return &identities.Identities, nil
}

func (c *Client) GetIdentity(ctx context.Context, local_part, id string) (*Identity, error) {
	var identity Identity

	body, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities/%s", local_part, id))
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	return &identity, nil
}

func (c *Client) CreateIdentity(ctx context.Context, local_part string, new_identity *Identity) (*Identity, error) {
	var identity Identity

	url_slug := fmt.Sprintf("mailboxes/%s/identities", local_part)

	json_body, err := json.Marshal(new_identity)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	body, err := c.Post(ctx, url_slug, json_body)
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
func (c *Client) CreateIdentityWithPassword(ctx context.Context, local_part, name, id, password string) (*Identity, error) {
	new_identity := Identity{Name: name, LocalPart: id}
	if password == "" {
		new_identity.Password = ""
		new_identity.PasswordUse = "mailbox"
	} else {
		new_identity.Password = password
		new_identity.PasswordUse = "custom"
	}

	return c.CreateIdentity(ctx, local_part, &new_identity)
}

// Convience function to create an identity that will not be used for authentication (i.e. login)
func (c *Client) CreateIdentityNoAuth(ctx context.Context, local_part, name, id string) (*Identity, error) {
	new_identity := Identity{Name: name, LocalPart: id, PasswordUse: "none"}

	return c.CreateIdentity(ctx, local_part, &new_identity)
}

func (c *Client) UpdateIdentity(ctx context.Context, local_part, id string, identity_params *Identity) (*Identity, error) {
	var identity Identity

	url_slug := fmt.Sprintf("mailboxes/%s/identities/%s", local_part, id)

	json_body, err := json.Marshal(identity_params)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	body, err := c.Put(ctx, url_slug, json_body)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	return &identity, nil
}

func (c *Client) DeleteIdentity(ctx context.Context, local_part, id string) error {
	url_slug := fmt.Sprintf("mailboxes/%s/identities/%s", local_part, id)

	_, err := c.Delete(ctx, url_slug)

	if err != nil {
		return fmt.Errorf("DeleteIdentity: %w", err)
	}

	return nil
}
