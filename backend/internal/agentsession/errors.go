package agentsession

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrSessionNotFound    = errx.New("agent session not found", errx.TypeNotFound)
	ErrSessionNotRunning  = errx.New("agent session not running", errx.TypeBusiness)
	ErrPresetNotInstalled = errx.New("preset not installed for tenant", errx.TypeBusiness)
)
