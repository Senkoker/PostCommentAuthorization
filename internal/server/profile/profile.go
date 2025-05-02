package profile

import (
	"VK_posts/internal/domain"
	"VK_posts/internal/models"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type PortProfile struct {
	PortProfileInterface PortProfileInterface
}
type PortProfileInterface interface {
	FillUserProfile(profile models.ProfileFill) (string, error)
	GetUserProfile(username string) (models.ProfileFill, error)
}

func NewProfileHandler(profileDomain *domain.ProfileDomain) *PortProfile {
	return &PortProfile{profileDomain}
}

func FillProfile(profilePort *PortProfile) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().ParseForm()
		_, img, err := c.Request().FormFile("img")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		birth := c.Request().FormValue("birth_date")
		profile := models.ProfileFill{
			UserID:     c.Get("userID").(string),
			FirstName:  c.Request().FormValue("first_name"),
			SecondName: c.Request().FormValue("second_name"),
			BirthDate:  birth,
			Education:  c.Request().FormValue("education"),
			Country:    c.Request().FormValue("country"),
			Image:      img,
			City:       c.Request().FormValue("city"),
		}
		fmt.Println(profile)
		id, err := profilePort.PortProfileInterface.FillUserProfile(profile)
		if err != nil {
			//Todo обработать ошибки
			return c.JSON(http.StatusInternalServerError, "")
		}
		return c.JSON(http.StatusOK, models.ProfileResponse{Id: id})

	}
}
func GetProfile(profilePort *PortProfile) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.QueryParam("userID")
		userProfile, err := profilePort.PortProfileInterface.GetUserProfile(userId)
		if err != nil {
			//todo: обработать ошибки
			return c.JSON(http.StatusInternalServerError, "")
		}
		return c.JSON(http.StatusOK, userProfile)
	}
}
