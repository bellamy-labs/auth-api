package v1

import (
	"github.com/bellamy-labs/auth-api/api/v1/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func (s *Server) initRoutes() {
	c := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}
	s.Handler.Use(middleware.CORSWithConfig(c))
	s.Handler.Use(middleware.Secure())
	// TODO: add CSRF protection
	// s.Handler.Use(middleware.CSRF())
	s.Handler.Use(middleware.BodyLimit("32M"))
	// TODO: implement Redis for in memory store - https://echo.labstack.com/middleware/rate-limiter/
	// TODO: add timeout functionality - https://echo.labstack.com/middleware/timeout/
	s.Handler.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	apiGroup := s.Handler.Group("/api/v1")

	authGroup := apiGroup.Group("/auth")

	googleGroup := authGroup.Group("/google")
	googleProvider, err := handlers.InitGoogle()
	if err != nil {
		log.Panic().Err(err)
	}
	googleGroup.GET("/", googleProvider.AuthHandler)
	googleGroup.POST("/", googleProvider.AuthHandler)
	googleGroup.GET("/callback", googleProvider.CallbackHandler)
	googleGroup.POST("/callback", googleProvider.CallbackHandler)
}
