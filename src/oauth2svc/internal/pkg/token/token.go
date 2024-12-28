package token

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/go-jose/go-jose/v4/jwt"
)

type IDToken struct {
	Iss        string           `json:"iss"`
	Sub        string           `json:"sub"`
	Name       string           `json:"name"`
	GivenName  string           `json:"given_name"`
	FamilyName string           `json:"family_name"`
	Email      string           `json:"email"`
	Aud        string           `json:"aud"`
	Exp        *jwt.NumericDate `json:"exp"`
	Iat        *jwt.NumericDate `json:"iat"`
}

type OpenIDConfiguration struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	UserinfoEndpoint       string   `json:"userinfo_endpoint"`
	RevocationEndpoint     string   `json:"revocation_endpoint"`
	JwksURI                string   `json:"jwks_uri"`
	ResponseTypesSupported []string `json:"response_types_supported"`
}

type MyClaims struct {
	JwtID  string           `json:"jti"`
	Scopes []string         `json:"scopes"`
	Roles  []string         `json:"roles"`
	Aud    string           `json:"aud"`
	Sub    string           `json:"sub"`
	Iss    string           `json:"iss"`
	Exp    *jwt.NumericDate `json:"exp"`
	Iat    *jwt.NumericDate `json:"iat"`
}

func New(l int) (t string, err error) {
	token := make([]byte, l)
	if _, err = rand.Read(token); err != nil {
		return
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
