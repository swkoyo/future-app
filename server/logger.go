package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func NewLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	zerolog := zerolog.New(output).With().Caller().Timestamp().Logger()
	Logger = zerolog
	return Logger
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		logger := Logger.With().Str("request_id", requestID).Logger()
		c.Set("logger", logger)

		logger.Info().Fields(map[string]interface{}{
			"method": c.Request().Method,
			"uri":    c.Request().URL.Path,
			"query":  c.Request().URL.RawQuery,
		}).Msg("Incoming request")

		err := next(c)
		if err != nil {
			logger.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")
			return err
		}

		return nil
	}
}

func GetEchoLogger(c echo.Context) zerolog.Logger {
	logger, _ := c.Get("logger").(zerolog.Logger)
	return logger
}
