package migagoapi

import (
	"encoding/json"
	"time"
)

const subpath = "mailboxes"

type CustomTime time.Time

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	*t = CustomTime(date)
	return
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"2006-01-02T15:04:05.000-0700"`)), nil
}

// TODO: if we're never marshaling into a mailbox struct,
// then 'omitempty' tag is really unnecessary because it
// does nothing when unmarshaling, only when marshaling is
// it relevant
type Mailbox struct {
	LocalPart             string    `json:"local_part,omitempty"`
	Domain                string    `json:"domain,omitempty"`
	Address               string    `json:"address,omitempty"`
	Name                  string    `json:"name,omitempty"`
	IsInternal            bool      `json:"is_internal,omitempty"`
	MaySend               bool      `json:"may_send,omitempty"`
	MayReceive            bool      `json:"may_receive,omitempty"`
	MayAccessImap         bool      `json:"may_access_imap,omitempty"`
	MayAccessPop3         bool      `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve  bool      `json:"may_access_managesieve,omitempty"`
	PasswordMethod        string    `json:"password_method,omitempty"`
	Password              string    `json:"password,omitempty"`
	PasswordRecoveryEmail string    `json:"password_recovery_email,omitempty"`
	SpamAction            string    `json:"spam_action,omitempty"`
	SpamAggressiveness    string    `json:"spam_aggressiveness"`
	SenderDenylist        string    `json:"sender_denylist,omitempty"`
	SenderAllowlist       string    `json:"sender_allowlist,omitempty"`
	RecipientDenylist     string    `json:"recipient_denylist,omitempty"`
	AutorespondActive     bool      `json:"autorespond_active,omitempty"`
	AutorespondSubject    string    `json:"autorespond_subject,omitempty"`
	AutorespondBody       string    `json:"autorespond_body,omitempty"`
	AutorespondExpiresOn  time.Time `json:"autorespond_expires_on,omitempty"`
	FooterActive          string    `json:"footer_active,omitempty"`
	FooterPlainBody       string    `json:"footer_plain_body,omitempty"`
	FooterHTMLBody        string    `json:"footer_html_body,omitempty"`
	Identities            []string  `json:"identities,omitempty"`
}
