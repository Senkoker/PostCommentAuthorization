package web

import (
	"VK_posts/internal/Postgres/handlers"
	Redis2 "VK_posts/internal/Redis"
	"VK_posts/internal/SelectelS3"
	"VK_posts/internal/config"
	"VK_posts/internal/domain"
	"VK_posts/internal/logger"
	server2 "VK_posts/internal/server"
	"VK_posts/internal/server/profile"
	"VK_posts/pkg/Postgres"
	"VK_posts/pkg/Redis"
)

func App() {
	cfg := config.NewConfig()
	logger.LoggerInit(cfg.Debug.DebugLogger)
	logger.GetLogger().Info("cfg info", cfg)

	//Storages
	postgresStorage := Postgres.NewStorage(cfg.PostgresUrl)
	redis := Redis.NewRedis(cfg)
	S3 := SelectelS3.NewSelectelS3(cfg)
	//PostCommentDomain
	postgresHandlers := handlers.NewStorageHandler(postgresStorage, cfg.PostgresCtx)
	redisHandler := Redis2.NewRedisHandler(redis, cfg.Redis.CtxTime)
	postsDomain := domain.NewDomain(postgresHandlers, postgresHandlers, S3, redisHandler)

	//ProfileDomain
	profilePostgresHandler := handlers.NewProfileHandler(postgresStorage)
	profileDomain := domain.NewProfileDomain(profilePostgresHandler, S3)
	profileHandler := profile.NewProfileHandler(profileDomain)

	server := server2.Server(cfg.Server.Host, cfg.Server.Port)
	server.FeedHandler(postsDomain)
	server.ProfileHandler(profileHandler)
	//server.MessengerHandler()
	server.Auth("sso_backend", 8080)
	server.Run()
}
