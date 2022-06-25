package v1

import (
	"github.com/bellamy-labs/auth-api/api/v1/handlers"
	"github.com/rs/zerolog/log"
)

func (s *Server) initRoutes() {
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
