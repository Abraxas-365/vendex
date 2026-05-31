package customergroup

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrGroupNotFound    = errx.New("customer group not found", errx.TypeNotFound)
	ErrAlreadyMember    = errx.New("customer is already a member of this group", errx.TypeConflict)
	ErrCustomerNotFound = errx.New("customer not found", errx.TypeNotFound)
	ErrMemberNotFound   = errx.New("group membership not found", errx.TypeNotFound)
)
