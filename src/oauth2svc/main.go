package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/smallbiznis/go-lib/pkg/errors"
	"github.com/smallbiznis/go-lib/pkg/server"
	v "github.com/smallbiznis/go-lib/pkg/validator"
	"github.com/smallbiznis/oauth2-server/internal/pkg/strings"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
	"github.com/smallbiznis/oauth2-server/service"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	env      string = "development"
	certFile string
	keyFile  string
)

func init() {
	env = strings.MustEnv("ENV", "development")
	certFile = strings.MustEnv("CERT_FILE", "/Users/taufiktriantono/go/src/microservice-demo/internal/smallbiz/cert.pem")
	keyFile = strings.MustEnv("KEY_FILE", "/Users/taufiktriantono/go/src/microservice-demo/internal/smallbiz/key.unencrypted.pem")
}

// Common modul
var CommonModule = fx.Options(
	fx.Provide(
		NewLogger,
		NewDatabase,
	),
	v.Validator,
	v.Translation,
	MigrationModule,
)

// Repository module
var RepositoryModule = fx.Options(
	fx.Provide(
		repository.NewApplicationRepository,
		repository.NewOrganizationRepository,
		repository.NewOrganizationKeyRepository,
		repository.NewAuthorizationCodeRepository,
		repository.NewAccessTokenRepository,
		repository.NewRefreshTokenRepository,
		repository.NewAccountRepository,
		repository.NewSessionRepository,
	),
)

// Service Module
var ServiceModule = fx.Module("service", fx.Options(
	fx.Provide(
		service.NewOAuthService,
		service.NewAccountService,
		service.NewApplicationService,
		service.NewOrganizationService,
	),
))

// Gin Module
var GinModule = fx.Module("gin", fx.Option(
	fx.Provide(
		NewGinEngine,
		NewOAuthHandler,
		NewAccountHandler,
		NewApplicationHandler,
	),
))

// Invoke Module
var InvokeModule = fx.Module("invoke", fx.Invoke(
	RegisterMiddleware,
	RegisterRoutes,
	StartServer,
))

var MigrationModule = fx.Module("migrate", fx.Invoke(
	StartAutoMigrate,
	SeedData,
))

func NewLogger() (log *zap.Logger) {
	fields := zap.Fields(
		zap.String("service_name", strings.MustEnv("SERVICE_NAME", "oauth2svc")),
		zap.String("service_version", strings.MustEnv("SERVICE_VERSION", "v1.0.0")),
		zap.String("service_environment", strings.MustEnv("ENV", "development")),
	)

	log, _ = zap.NewProduction(fields)
	if env == "development" {
		log = zap.NewExample(fields)
	}

	return
}

func NewDatabase(log *zap.Logger) (db *gorm.DB, err error) {
	dsn := postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			strings.MustEnv("DB_HOST", "127.0.0.1"),
			strings.MustEnv("DB_USER", "postgres"),
			strings.MustEnv("DB_PASSWORD", "35411231"),
			strings.MustEnv("DB_DB", "oauth2"),
			strings.MustEnv("DB_PORT", "5432"),
			strings.MustEnv("DB_SSL_MODE", "disable"),
			strings.MustEnv("DB_TIMEZONE", "Asia/Jakarta"),
		))

	db, err = gorm.Open(dsn)
	if err != nil {
		zap.Error(err)
		return
	}

	if env == "development" {
		db = db.Debug()
	}

	return
}

func NewGinEngine() (*gin.Engine, http.Handler) {
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	return app, app.Handler()
}

func RegisterRoutes(
	r *gin.Engine,
	s *OAuthHandler,
	a *AccountHandler,
	app *ApplicationHandler,
) {
	wellKnow := r.Group("/.well-known")
	{
		wellKnow.GET("/openid-configuration", s.HandleOpenIDConfiguration)
		wellKnow.GET("/jwks.json", s.HandleGetKeys)
	}

	oauth := r.Group("/oauth")
	{
		oauth.GET("/introspect", s.HandleIntrospect)
		oauth.GET("/revoke", s.HandleRevoke)
		oauth.GET("/userinfo", s.HandleUserInfo)
		oauth.POST("/token", s.HandleRequestToken)
		oauth.GET("/authorize", s.HandleRequestAuthorization)
	}

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			applications := v1.Group("/applications")
			{
				applications.GET("", app.HandleList)
				applications.POST("", app.HandleCreate)
				applications.GET("/:appId", app.HandleGet)
				applications.PUT("/:appId", app.HandleUpdate)
				applications.DELETE("/:appId", app.HandleDelete)
			}

			password := v1.Group("/password")
			{
				password.GET("/rules")
			}

			accounts := v1.Group("/accounts")
			{
				accounts.GET("/lookup", a.HandleLookup)
				accounts.POST("/signup", a.HandleSignUp)
				accounts.POST("/signInWithPassword", a.HandleSignInWithPassword)
				accounts.POST("/sendVerificationCode", a.HandleSendVerificationCode)
				accounts.POST("/signInWithPhoneNumber", a.HandleSignInWithPhoneNumber)
			}
		}
	}
}

func RegisterMiddleware(organizationRepository repository.IOrganizationRepository, r *gin.Engine, log *zap.Logger, translate ut.Translator) {
	r.Use(RegisterMiddlewareTenant(organizationRepository))
	// r.Use(middleware.GinZapLogger(log))
	r.Use(MiddlewareError(translate))
}

func RegisterMiddlewareTenant(organizationRepository repository.IOrganizationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		host := c.Request.Host
		subdomain := strings.Split(host, ".")

		orgID := strings.TrimSpace(subdomain[0])
		exist, err := organizationRepository.FindOne(ctx, model.Organization{
			Name: orgID,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}

		if exist == nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(ctx, "tenant", exist))

		c.Next()
	}
}

func MiddlewareError(translate ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			c.JSON(
				ValidationError(err.Err, translate),
			)
		}
	}
}

func ValidationError(err error, translate ut.Translator) (code int, obj any) {
	code = 500
	obj = gin.H{
		"error": errors.InternalServerError("InternalServerError", err.Error()),
	}

	// Handle error io.EOF request body empty
	if err == io.EOF {
		code = 400
		obj = gin.H{
			"error": errors.BadRequest("InvalidRequest", "request can't be empty"),
		}
		return
	}

	// Handle error *json.SyntaxError
	if _, ok := err.(*json.SyntaxError); ok {
		code = 400
		obj = gin.H{
			"error": errors.BadRequest("InvalidRequest", err.Error()),
		}
		return
	}

	// Handle error *validator.fieldError
	if e, ok := err.(validator.FieldError); ok {
		code = 400
		msg := fmt.Errorf(e.Translate(translate)).Error()
		// newErr := errors.BadRequest("InvalidRequest", msg)
		obj = gin.H{
			"error": gin.H{
				"status":  code,
				"name":    "InvalidRequest",
				"message": msg,
				"details": []gin.H{
					{
						"field": e.Field(),
						"tags":  e.Tag(),
					},
				},
			},
		}
		return
	}

	// Handle error uuid.isInvalidLength
	if uuid.IsInvalidLengthError(err) {
		code = 400
		obj = gin.H{
			"error": err,
		}
		return
	}

	// Handle error business logic
	if _, ok := err.(errors.Error); ok {
		code = 400
		obj = gin.H{
			"error": err,
		}
		return
	}

	return
}

func StartAutoMigrate(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(
		&model.Provider{},
		&model.Organization{},
		&model.OrganizationKey{},
		&model.Account{},
		&model.Application{},
		&model.Profile{},
		&model.AccessToken{},
		&model.RefreshToken{},
		&model.AuthorizationCode{},
		&model.UserSession{},
	); err != nil {
		if err.Error() != "insufficient arguments" {
			return
		}

		return nil
	}

	return nil
}

func StartServer(log *zap.Logger, s server.IServer, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			log.Info("starting application")
			err = s.Run()
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			log.Info("shutdown application")
			err = s.Down(ctx)
			return
		},
	})
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	opts := []fx.Option{
		CommonModule,
		RepositoryModule,
		ServiceModule,
		GinModule,
		InvokeModule,
		server.Module,
	}

	if env == "development" {
		opts = append(opts, fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}))
	}

	app := fx.New(opts...)
	go func() {
		if err := app.Start(ctx); err != nil {
			zap.Error(err)
		}
	}()

	<-ctx.Done()
	if err := app.Stop(context.TODO()); err != nil {
		zap.Error(err)
	}
}
