package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	internalErr "github.com/smallbiznis/go-lib/pkg/errors"
	"github.com/smallbiznis/oauth2-server/internal/pkg/errors"
	"github.com/smallbiznis/oauth2-server/internal/pkg/strings"
	"github.com/smallbiznis/oauth2-server/internal/pkg/token"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
)

type IOAuthService interface {
	HandleIntrospect(c *gin.Context)
	HandleUserInfo(c *gin.Context)
	HandleRevoke(c *gin.Context)
	HandleRequestAuthorization(c *gin.Context)
	HandleRequestToken(ctx context.Context, req model.TokenRequest) (resp gin.H, err error)
	HandleGetKeys(ctx context.Context) ([]jose.JSONWebKey, error)
	HandleOpenIDConfiguration(c *gin.Context)
}

type oauthService struct {
	validator                   *validator.Validate
	organizationKey             repository.IOrganizationKeyRepository
	applicationRepository       repository.IApplicationRepository
	authorizationCodeRepository repository.IAuthorizationCodeRepository
	accessTokenRepository       repository.IAccessTokenRepository
	refreshTokenRepository      repository.IRefreshTokenRepository
	accountRepository           repository.IAccountRepository
	sessionRepository           repository.ISessionRepository
}

func NewOAuthService(
	validator *validator.Validate,
	organizationKeyRepository repository.IOrganizationKeyRepository,
	applicationRepository repository.IApplicationRepository,
	authorizationCodeRepository repository.IAuthorizationCodeRepository,
	accessTokenRepository repository.IAccessTokenRepository,
	refreshTokenRepository repository.IRefreshTokenRepository,
	accountRepository repository.IAccountRepository,
	sessionRepository repository.ISessionRepository,
) IOAuthService {
	return &oauthService{
		validator,
		organizationKeyRepository,
		applicationRepository,
		authorizationCodeRepository,
		accessTokenRepository,
		refreshTokenRepository,
		accountRepository,
		sessionRepository,
	}
}

func (svc *oauthService) HandleIntrospect(c *gin.Context) {}

func (svc *oauthService) HandleUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	split := strings.Split(header, " ")
	if len(split) < 1 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if split[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := svc.accessTokenRepository.FindOne(ctx, model.AccessToken{OrganizationID: tenant.ID, AccessToken: split[1]})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if t == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := svc.accountRepository.FindOne(ctx, model.Account{ID: t.UserID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (svc *oauthService) HandleRevoke(c *gin.Context) {}

func (svc *oauthService) HandleRequestAuthorization(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	var req model.AuthorizationRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	if err := svc.validator.Struct(&req); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	query, _ := query.Values(&req)
	exist, err := svc.applicationRepository.FindOne(ctx, []string{}, model.Application{
		OrganizationID: tenant.ID,
		ClientID:       req.ClientID,
	})
	if err != nil {
		// handle error
		c.Error(err)
		return
	}

	if exist == nil {
		// handle error application not found
		c.Error(errors.ErrInvalidCredential)
		return
	}

	if len(exist.RedirectUrls) == 0 {
		// handle error misconfiguration
		c.Error(errors.ErrMissingConfiguration)
		return
	}

	redirectUris := make(map[string]bool, 0)
	for _, v := range exist.RedirectUrls {
		if v == req.RedirectUri {
			redirectUris[v] = true
		} else {
			redirectUris[v] = false
		}
	}

	// handle invalid redirect_uri
	if !redirectUris[req.RedirectUri] {
		c.Error(errors.ErrInvalidRedirectURI)
		return
	}

	scopes := strings.Split(req.Scope, " ")
	v, err := c.Cookie("_SID")
	if v == "" || err != nil && err != http.ErrNoCookie {
		c.Header("Location", fmt.Sprintf("/signin?%s", query.Encode()))
		c.Status(http.StatusFound)
		return
	}

	sess, err := svc.sessionRepository.FindOne(ctx, model.UserSession{ID: v})
	if err != nil {
		c.Header("Location", fmt.Sprintf("/signin?%s", query.Encode()))
		c.Status(http.StatusFound)
		return
	}

	if sess == nil {
		c.Header("Location", fmt.Sprintf("/signin?%s", query.Encode()))
		c.Status(http.StatusFound)
		return
	}

	authCode := model.NewAuthorizationCode()
	authCode.OrganizationID = exist.OrganizationID
	authCode.ApplicationID = exist.ID
	authCode.UserID = sess.UserID
	authCode.Scopes = scopes
	authCode.Code = strings.RandomHex(16)
	authCode.RedirectUri = req.RedirectUri
	authCode.CodeChallenge = req.CodeChallenge
	authCode.CodeChallengeMethod = req.CodeChallengeMethod

	authorize, err := svc.authorizationCodeRepository.Save(ctx, *authCode)
	if err != nil {
		// handle error
		c.Error(err)
		return
	}

	url := make(url.Values)
	url.Set("code", authorize.Code)
	if req.State != "" {
		url.Set("state", req.State)
	}

	c.Header("Location", fmt.Sprintf("%s?%s", authorize.RedirectUri, url.Encode()))
	c.Status(http.StatusFound)
}

func (svc *oauthService) HandleRequestToken(ctx context.Context, req model.TokenRequest) (resp gin.H, err error) {
	tenant := ctx.Value("tenant").(*model.Organization)

	app, err := svc.applicationRepository.FindOne(ctx, []string{}, model.Application{
		OrganizationID: tenant.ID,
		ClientID:       req.ClientID,
	})
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, errors.ErrInvalidCredential
	}

	if app.ClientSecret != req.ClientSecret {
		return nil, errors.ErrInvalidCredential
	}

	tgr := model.TokenGeneration{
		Organization: tenant.ID,
		Application:  app.Name,
	}

	switch req.GrantType {
	case model.GrantAuthorizationCode:
		if req.Code == "" {
			return nil, errors.ErrInvalidAuthorizationCode
		}

		authCode, err := svc.authorizationCodeRepository.FindOne(ctx, model.AuthorizationCode{Code: req.Code})
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		if authCode == nil {
			return nil, errors.ErrInvalidAuthorizationCode
		}

		if authCode.Revoke {
			return nil, internalErr.BadRequest("InvalidRequest", "authorization code has been revoked")
		}

		user, err := svc.accountRepository.FindOne(ctx, model.Account{
			ID: authCode.UserID,
		})
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		if user == nil {
			return nil, internalErr.BadRequest("InvalidRequest", "invalid authorization code")
		}

		authCode.Revoke = !authCode.Revoke
		if _, err := svc.authorizationCodeRepository.Update(ctx, *authCode); err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		tgr.User = *user
		tgr.Scope = authCode.Scopes

	case model.GrantClientCredentials:
		user, err := svc.accountRepository.FindOne(ctx, model.Account{
			ID:             *app.UserID,
			OrganizationID: app.OrganizationID,
		})
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		if user == nil {
			return nil, internalErr.BadRequest("InvalidRequest", "invalid credential")
		}

		tgr.User = *user
		for _, v := range strings.Split(strings.TrimSpace(req.Scope), " ") {
			tgr.Scope = append(tgr.Scope, v)
		}

	case model.GrantRefreshToken:
		refreshToken, err := svc.refreshTokenRepository.FindOne(ctx, model.RefreshToken{
			RefreshToken: &req.RefreshToken,
		})
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		if refreshToken == nil {
			return nil, internalErr.BadRequest("InvalidRequest", "invalid refresh token")
		}

		if refreshToken.Revoke {
			return nil, internalErr.BadRequest("InvalidRequest", "refresh token has beed revoked")
		}

		if refreshToken.RefreshTokenExpiresIn.Compare(time.Now()) > -1 {
			return nil, internalErr.BadRequest("InvalidRequest", "refresh token has beend expired")
		}

		refreshToken.Revoke = !refreshToken.Revoke
		refreshToken.UpdatedAt = time.Now()

		if _, err := svc.refreshTokenRepository.Update(ctx, *refreshToken); err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		tgr.User = model.AggregateAccount{Account: model.Account{ID: refreshToken.UserID}}
		tgr.Scope = refreshToken.Scopes

	case model.GrantPassword:

		account, err := svc.accountRepository.FindOne(ctx, model.Account{
			OrganizationID: tgr.Organization,
			Provider:       model.Password,
			Username:       req.Username,
		})
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		if account == nil {
			return nil, internalErr.BadRequest("InvalidRequest", "invalid email or password")
		}

		if !account.Account.ComparePassword(req.Password) {
			return nil, internalErr.BadRequest("InvalidRequest", "invalid email or password")
		}

		tgr.User = *account
		tgr.Scope = account.Roles

	default:
		return nil, internalErr.BadRequest("UnsupportedGrantType", "unsupported grant type")
	}

	keys, err := svc.organizationKey.FindOne(ctx, model.OrganizationKey{
		OrganizationID: tenant.ID,
		Use:            model.Sign,
	})
	if err != nil {
		return nil, internalErr.InternalServerError("InternalServerError", err.Error())
	}

	if keys == nil {
		return nil, internalErr.InternalServerError("InternalServerError", "key not found")
	}

	jwtId := uuid.NewString()
	expiresIn := time.Now().Add(time.Second * time.Duration(app.AccessTokenExpiresIn))
	privateKey, err := svc.parseRSAPrivateKeyFromPEM(keys.Key)
	if err != nil {
		return nil, err
	}

	signKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       privateKey,
	}

	issuer := fmt.Sprintf("https://%s.smallbiznis.test", app.Name)
	idTokenClaims := token.IDToken{
		Iss:        issuer,
		Sub:        tgr.User.ID,
		Name:       fmt.Sprintf("%s %s", tgr.User.FirstName, tgr.User.LastName),
		GivenName:  tgr.User.FirstName,
		FamilyName: tgr.User.LastName,
		Email:      tgr.User.Username,
		Aud:        app.ClientID,
		Exp:        jwt.NewNumericDate(expiresIn),
		Iat:        jwt.NewNumericDate(time.Now()),
	}

	idToken, err := svc.generateToken(signKey, keys.KeyID, idTokenClaims)
	if err != nil {
		return nil, internalErr.InternalServerError("InternalServerError", err.Error())
	}

	claims := token.MyClaims{
		Iss:    issuer,
		Scopes: tgr.Scope,
		JwtID:  jwtId,
		Aud:    app.ClientID,
		Sub:    tgr.User.ID,
		Roles:  tgr.User.Roles,
		Exp:    jwt.NewNumericDate(expiresIn),
		Iat:    jwt.NewNumericDate(time.Now()),
	}

	t, err := svc.generateToken(signKey, keys.KeyID, claims)
	if err != nil {
		return nil, internalErr.InternalServerError("InternalServerError", err.Error())
	}

	var (
		scopes    string
		isOffline bool
	)
	for i, v := range tgr.Scope {
		if v == "offline_access" {
			isOffline = true
		}

		if i != 0 {
			scopes += fmt.Sprintf(" %s", v)
		} else {
			scopes += v
		}
	}

	accessToken := model.AccessToken{
		OrganizationID:       tenant.ID,
		ApplicationID:        app.ID,
		UserID:               tgr.User.ID,
		Scopes:               tgr.Scope,
		AccessToken:          t,
		AccessTokenExpiresIn: expiresIn,
	}

	resp = gin.H{
		"id_token":     idToken,
		"token_type":   "Bearer",
		"access_token": t,
		"expires":      app.AccessTokenExpiresIn,
	}

	if isOffline {
		expiresIn := time.Now().Add(7 * 24 * time.Hour)
		t, err := token.New(32)
		if err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		refreshToken := model.RefreshToken{
			OrganizationID:        tenant.ID,
			ApplicationID:         app.ID,
			UserID:                tgr.User.ID,
			Scopes:                tgr.Scope,
			RefreshToken:          &t,
			RefreshTokenExpiresIn: &expiresIn,
		}

		if _, err := svc.refreshTokenRepository.Save(ctx, refreshToken); err != nil {
			return nil, internalErr.InternalServerError("InternalServerError", err.Error())
		}

		resp["refresh_token"] = refreshToken
	}

	if _, err := svc.accessTokenRepository.Save(ctx, accessToken); err != nil {
		return nil, internalErr.InternalServerError("InternalServerError", err.Error())
	}

	resp["scope"] = scopes
	return resp, nil
}

func (svc *oauthService) HandleRequestAuthorizationCode(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(model.Organization)

	sess := sessions.Default(c)
	v := sess.Get("user")

	req := model.AuthorizationRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	app, err := svc.applicationRepository.FindOne(ctx, []string{}, model.Application{
		OrganizationID: tenant.ID,
		ClientID:       req.ClientID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, internalErr.InternalServerError("InternalServerError", err.Error()))
		return
	}

	if app == nil {
		c.JSON(http.StatusBadRequest, internalErr.BadRequest("InvalidClient", "invalid client"))
		return
	}

	redirectUris := make(map[string]bool, 0)
	for _, v := range app.RedirectUrls {
		if v == req.RedirectUri {
			redirectUris[v] = true
		} else {
			redirectUris[v] = false
		}
	}

	if !redirectUris[req.RedirectUri] {
		c.JSON(http.StatusBadRequest, internalErr.BadRequest("UnauthorizedClient", "callback_url missmatch"))
		return
	}

	scopes := strings.Split(req.Scope, " ")
	query, _ := query.Values(&req)
	if v == nil {
		sess.Set("state", query.Encode())
		sess.Save()
		c.JSON(http.StatusUnauthorized, internalErr.Unauthorized("Unauthorized", "Your are not authenticated!"))
		return
	}

	authCode := model.AuthorizationCode{
		OrganizationID:      app.OrganizationID,
		ApplicationID:       app.Name,
		UserID:              v.(string),
		Scopes:              scopes,
		Code:                strings.RandomHex(16),
		RedirectUri:         req.RedirectUri,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	url := url.Values{}
	url.Set("code", authCode.Code)
	if req.State != "" {
		url.Set("state", req.State)
	}

	c.Header("Location", fmt.Sprintf("%s?%s", authCode.RedirectUri, url.Encode()))
	c.Status(http.StatusOK)
}

func (svc *oauthService) HandleOpenIDConfiguration(c *gin.Context) {
	req := c.Request
	issuer := fmt.Sprintf("https://%s", req.Host)
	c.JSON(http.StatusOK, token.OpenIDConfiguration{
		Issuer:                 issuer,
		AuthorizationEndpoint:  fmt.Sprintf("%s/oauth/authorize", issuer),
		TokenEndpoint:          fmt.Sprintf("%s/oauth/token", issuer),
		UserinfoEndpoint:       fmt.Sprintf("%s/oauth/userinfo", issuer),
		RevocationEndpoint:     fmt.Sprintf("%s/oauth/revoke", issuer),
		JwksURI:                fmt.Sprintf("%s/.well-known/jwks.json", issuer),
		ResponseTypesSupported: []string{"code"},
	})
}

func (svc *oauthService) parseRSAPrivateKeyFromPEM(privateKey string) (*rsa.PrivateKey, error) {
	// Decode private key dari PEM format
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, internalErr.InternalServerError("FailedParseRSAFromPEM", "failed to decode PEM block containing private key")
	}

	pv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Jika gagal dengan PKCS1, coba dengan PKCS8
		privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, internalErr.InternalServerError("FailedParseRSAFromPEM", fmt.Sprintf("failed to parse RSA private key: %v", err))
		}

		// Casting dari interface{} ke *rsa.PrivateKey
		pv, ok := privateKeyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, internalErr.InternalServerError("FailedParseRSAFromPEM", "not an RSA private key")
		}

		return pv, nil
	}

	return pv, nil
}

func (svc *oauthService) generateToken(signKey jose.SigningKey, keyId string, claims interface{}) (t string, err error) {

	opt := jose.SignerOptions{}
	opt.WithHeader(jose.HeaderKey("kid"), keyId)

	sign, err := jose.NewSigner(signKey, &opt)
	if err != nil {
		fmt.Printf("failed jose.NewSigner: %v\n", err)
		return "", err
	}

	builder := jwt.Signed(sign)
	builder = builder.Claims(claims)
	return builder.Serialize()
}

func (svc *oauthService) HandleGetKeys(ctx context.Context) (jsonwebKeys []jose.JSONWebKey, err error) {
	tenant := ctx.Value("tenant").(*model.Organization)
	jwks, _, err := svc.organizationKey.Find(ctx, nil, model.OrganizationKey{
		OrganizationID: tenant.ID,
		Use:            model.Sign,
	})
	if err != nil {
		return nil, internalErr.InternalServerError("InternalServerError", err.Error())
	}

	var joseJwks []jose.JSONWebKey
	for _, v := range jwks {

		privateKey, err := svc.parseRSAPrivateKeyFromPEM(v.Key)
		if err != nil {
			return jsonwebKeys, err
		}

		j := jose.JSONWebKey{
			Algorithm: string(v.Algorithm),
			Key:       &privateKey.PublicKey,
			KeyID:     v.KeyID,
			Use:       v.Use.String(),
		}

		if _, err := j.MarshalJSON(); err != nil {
			fmt.Printf("failed: MarshalJSON %v\n", err)
		}

		joseJwks = append(joseJwks, j)
	}

	return joseJwks, nil
}
