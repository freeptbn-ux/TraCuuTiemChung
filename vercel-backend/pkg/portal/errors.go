package portal

import "errors"

var (
	ErrLoginFailed     = errors.New("portal login failed")
	ErrSessionExpired  = errors.New("portal session expired")
	ErrPortalDown      = errors.New("portal is currently unavailable")
	ErrPatientNotFound = errors.New("patient not found")
)
