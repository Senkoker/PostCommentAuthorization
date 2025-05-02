package middlewares

import (
	"VK_posts/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

const (
	secret = "my_secret"
)

func JWTParser(token string) (string, error) {
	parseToken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid signature token")
		}
		logger.GetLogger().Info("Token parse error", "err", err)
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Token parse error")
	}
	claims := parseToken.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))

	userID := claims["user.id"].(string)
	timeNow := time.Now().Unix()
	if timeNow > exp {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Token expired")
	}
	return userID, nil
}

func InformationAboutRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		url := c.Request().URL.String()
		ip := c.Request().RemoteAddr
		timeStart := time.Now()
		logger.GetLogger().Info("Data of request", "method", method, "ip", ip, "url", url)
		err := next(c)
		if err != nil {
			logger.GetLogger().Error("Problem with request", "method", method, "ip", ip, "url", url, "err", err)
		}
		logger.GetLogger().Info("RequestInfo", "method", method, "ip", ip, "url", url)
		timeNow := time.Now()
		logger.GetLogger().Info("Time of request", "Interval", timeNow.Sub(timeStart).Seconds())
		return err
	}
}
func CheckTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenBearer := c.Request().Header.Get("Authorization")
		token := tokenBearer[7:]
		userId, err := JWTParser(token)
		if err != nil {
			return err
		}
		c.Set("userID", userId)
		return next(c)

	}
}
