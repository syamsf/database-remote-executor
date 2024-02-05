package config

import (
	"database-remote-commander/internal/delivery/http"
	"database-remote-commander/internal/delivery/http/middleware"
	"database-remote-commander/internal/delivery/http/route"
	"database-remote-commander/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type BootstrapConfig struct {
	App    *fiber.App
	Log    *logrus.Logger
	Config *viper.Viper
	DB     *mongo.Database
}

func Bootstrap(config *BootstrapConfig) {
	// Setup middleware
	headerAuthMiddleware := middleware.NewAuth(config.Log, config.Config)

	// setup use cases
	remoteQueryUseCase := usecase.NewRemoteQueryUseCase(config.DB, config.Log, config.Config)

	// setup controller
	remoteQueryController := http.NewQueryController(remoteQueryUseCase, config.Log)

	routeConfig := route.Config{
		App:                  config.App,
		HeaderAuthMiddleware: headerAuthMiddleware,
		QueryController:      remoteQueryController,
	}

	routeConfig.Setup()
}
