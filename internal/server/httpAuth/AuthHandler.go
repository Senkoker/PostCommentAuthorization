package httpAuth

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"context"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"github.com/labstack/echo"
	"net/http"
)

func LoginHandler(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		var login models.LoginInfo
		err := c.Bind(&login)
		logger.GetLogger().Error("eerr", err)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		loginResp, err := client.NewAuthClient.Login(ctx, &sso_v1_ssov1.Loginrequest{Email: login.Email, Password: login.Password, Appid: login.App})
		if err != nil {
			logger.GetLogger().Info("Login Error", "err", err.Error())
			return c.JSON(http.StatusOK, models.ResponseError{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, models.LoginResponse{AccessToken: loginResp.Token})
	}
}

func RegisterHandler(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		var register models.RegisterInfo
		err := c.Bind(&register)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		loginResp, err := client.NewAuthClient.Register(ctx, &sso_v1_ssov1.Registrequest{Email: register.Email, Password: register.Password})
		if err != nil {
			//todo: Обработать ошибки как в GRPC
			logger.GetLogger().Info("Login Error", "err", err.Error())
			return c.JSON(http.StatusOK, models.ResponseError{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, models.RegisterResponse{Id: loginResp.Userid})
	}

}
func Accept(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		var accept models.AcceptInfo
		err := c.Bind(&accept)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		acceptResp, err := client.NewAuthClient.Accept(ctx, &sso_v1_ssov1.Acceptrequest{Usercode: accept.UserCode})
		if err != nil {
			logger.GetLogger().Info("Login Error", "err", err.Error())
			return c.JSON(http.StatusOK, models.ResponseError{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, models.AcceptResponse{Status: acceptResp.Accresp})
	}
}
