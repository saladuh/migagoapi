package migagoapi

type Addresser interface {
	GetAddress() string
}

func (m *Mailbox) GetAddress() string {
	return m.Address
}

func (i *Identity) GetAddress() string {
	return i.Address
}
