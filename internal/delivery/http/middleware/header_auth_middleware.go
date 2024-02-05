package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewAuth(logger *logrus.Logger, config *viper.Viper) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headerKey := ctx.Get("X-REMOTE-COMMANDER-AUTH", "NOT_FOUND")
		validationHeader := config.GetString("auth.header_key")
		logger.Debugf("Header X-REMOTE-COMMANDER-AUTH: %s", headerKey)

		if headerKey != validationHeader {
			logger.Warnf("Header validation failed: %s", ctx.IP())
			return fiber.ErrUnauthorized
		}

		logger.Debugf("Header validation success: %s", ctx.IP())

		return ctx.Next()
	}
}
