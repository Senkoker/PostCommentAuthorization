package httpAuth

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"context"
	"encoding/json"
	"fmt"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func LoginHandler(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		var login models.LoginInfo
		err := c.Bind(&login)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		loginResp, err := client.NewAuthClient.Login(ctx, &sso_v1_ssov1.Loginrequest{Email: login.Email, Password: login.Password, Appid: login.App})
		if err != nil {
			logger.GetLogger().Info("Login Error", "err", err.Error())
			if strings.Contains(err.Error(), "Email or Password is empty") {
				return echo.NewHTTPError(http.StatusBadRequest, "Email or Password is empty")
			}
			if strings.Contains(err.Error(), "This user is not exist") {
				return echo.NewHTTPError(http.StatusBadRequest, "This user is not exist")
			}
			if strings.Contains(err.Error(), "Invalid email or password") {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid email or password")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Problem in server")
		}
		return c.JSON(http.StatusOK, models.LoginResponse{AccessToken: loginResp.Token})
	}
}

func RegisterHandler(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		var register models.RegisterInfo
		err := json.NewDecoder(c.Request().Body).Decode(&register)
		fmt.Println(register)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		loginResp, err := client.NewAuthClient.Register(ctx, &sso_v1_ssov1.Registrequest{Email: register.Email, Password: register.Password})
		if err != nil {
			logger.GetLogger().Info("Login Error", "err", err.Error())
			if strings.Contains(err.Error(), "Password or login is empty") {
				return echo.NewHTTPError(http.StatusBadRequest, "Email or Password is empty")
			}
			if strings.Contains(err.Error(), "This user already get code in his email") {
				return echo.NewHTTPError(http.StatusConflict, "This user already get code in his email")
			}
			if strings.Contains(err.Error(), "This user already exists") {
				return echo.NewHTTPError(http.StatusConflict, "This user already exists")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Problem in server")
		}
		return c.JSON(http.StatusOK, models.RegisterResponse{Id: loginResp.Userid})
	}

}
func Accept(client *AuthGrpcConnect) echo.HandlerFunc {
	return func(c echo.Context) error {
		userCode := c.QueryParam("accept")
		if userCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		ctx := context.Background()
		_, err := client.NewAuthClient.Accept(ctx, &sso_v1_ssov1.Acceptrequest{Usercode: userCode})
		if err != nil {
			logger.GetLogger().Info("Login Error", "err", err.Error())
			if strings.Contains(err.Error(), "This user already accept his data") {
				return echo.NewHTTPError(http.StatusConflict, "This user already accept his data")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Problem in server")
		}
		return c.HTML(200, "<h1>Регистрация подтверждена можете зайти в свою учетную запись</h1>")
	}
}
