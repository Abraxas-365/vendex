package kernel

import (
	"fmt"
	"net/mail"
)

type Email string

func NewEmail(raw string) (Email, error) {
	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return "", fmt.Errorf("invalid email: %q", raw)
	}
	return Email(addr.Address), nil
}

func (e Email) String() string { return string(e) }
