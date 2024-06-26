package api

import (
	_ "api-gateway/api/docs" // swag
	"api-gateway/api/handlers/middleware"
	"api-gateway/api/handlers/tokens"
	v1 "api-gateway/api/handlers/v1"
	"api-gateway/config"
	"api-gateway/pkg/logger"
	"api-gateway/services"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	CasbinEnforcer *casbin.Enforcer
}

// @Title Welcome to swagger service
// @Version 1.0
// @Description you can use this as social network
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	jwtHandler := tokens.JWTHandler{
		SigninKey: option.Conf.SigningKey,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		JWTHandler:     jwtHandler,
		Enforcer:       option.CasbinEnforcer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.NewAuthorizer(option.CasbinEnforcer, jwtHandler, option.Conf))

	api := router.Group("/v1")

	// users
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)
	api.GET("/users", handlerV1.ListUsers)
	api.PUT("/users/:id", handlerV1.UpdateUser)
	api.DELETE("/users/:id", handlerV1.DeleteUser)

	// user registratsiya
	api.POST("/register", handlerV1.Register)
	api.GET("/login", handlerV1.LogIn)
	api.GET("/verification", handlerV1.Verification)
	api.GET("/refreshusertoken", handlerV1.RefreshUserToken)

	url := ginSwagger.URL("swaggerdoc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
