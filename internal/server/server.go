package server

import (
	"VK_posts/internal/domain"
	"VK_posts/internal/server/httpAuth"
	"VK_posts/internal/server/httpPostAndComment"
	"VK_posts/internal/server/middlewares"
	profile2 "VK_posts/internal/server/profile"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Router struct {
	router  *echo.Echo
	address string
}

func Server(host string, port int) *Router {
	e := echo.New()
	address := net.JoinHostPort(host, strconv.Itoa(port))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	return &Router{e, address}
}
func (e *Router) Run() {
	err := e.router.Start(e.address)
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Router) FeedHandler(postDomain *domain.Domain) {
	postHandler := httpPostAndComment.NewHandlers(postDomain)
	feed := e.router.Group("/feed")
	feed.Use(middlewares.InformationAboutRequest)
	feed.Use(middlewares.CheckTokenMiddleware)
	feed.POST("/create_post", postHandler.CreatePost())
	feed.POST("/create_comment", postHandler.CreateComment())
	feed.POST("/get_posts", postHandler.GetPosts())
	feed.GET("/get_main_comments", postHandler.GetMainComment())
	feed.GET("/get_comments", postHandler.GetComment())
}

func (e *Router) Auth(address string, port int) {
	authClient := httpAuth.NewAuthClient(address, port)
	auth := e.router.Group("/auth")
	auth.Use(middlewares.InformationAboutRequest)
	auth.POST("/login", httpAuth.LoginHandler(authClient))
	auth.POST("/accepter", httpAuth.Accept(authClient))
	auth.POST("/register", httpAuth.RegisterHandler(authClient))
}
func (e *Router) ProfileHandler(profilePort *profile2.PortProfile) {
	profile := e.router.Group("/profile")
	profile.Use(middlewares.CheckTokenMiddleware)
	profile.POST("/fill_data", profile2.FillProfile(profilePort))
	profile.GET("/get_profile", profile2.GetProfile(profilePort))
}

func (e *Router) MessengerHandler() {}
