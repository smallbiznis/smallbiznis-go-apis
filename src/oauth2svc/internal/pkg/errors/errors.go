package errors

import "github.com/smallbiznis/go-lib/pkg/errors"

var (
	notImplement               = "NotImplement"
	invalidCredential          = "InvalidCredential"
	invalidRedirectURI         = "InvalidRedirectURI"
	missingConfiguration       = "MissingConfiguration"
	invalidAuthorizationCode   = "InvalidAuthorizationCode"
	invalidEmailOrPassword     = "InvalidEmailOrPassword"
	invalidVerificationCode    = "InvalidVerificationCode"
	verificationCodeExpired    = "VerificationCodeExpired"
	invalidApplicationNotFound = "InvalidApplicationNotFound"
)

var (
	ErrNotImplement             = errors.BadRequest(notImplement, "Not Implement")
	ErrInvalidCredential        = errors.BadRequest(invalidCredential, "invalid credential")
	ErrMissingConfiguration     = errors.BadRequest(missingConfiguration, "missing configuration")
	ErrInvalidRedirectURI       = errors.BadRequest(invalidRedirectURI, "invalid `redirect_uri`")
	ErrInvalidAuthorizationCode = errors.BadRequest(invalidAuthorizationCode, "invalid authorization code")
	ErrInvalidEmailOrPassord    = errors.BadRequest(invalidEmailOrPassword, "invalid email or password")
	ErrInvalidVerificationCode  = errors.BadRequest(invalidVerificationCode, "invalid verification code")
	ErrVerificationCodeExpired  = errors.BadRequest(verificationCodeExpired, "expired session_id")
	ErrApplicationNotFound      = errors.BadRequest(invalidApplicationNotFound, "application not found")
)
