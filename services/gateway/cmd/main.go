package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/zvynt/gateway/config"
	"github.com/junaid9001/zvynt/gateway/handlers"
	"github.com/junaid9001/zvynt/gateway/middlewares"
	redisclient "github.com/junaid9001/zvynt/gateway/redis-client"
	"github.com/junaid9001/zvynt/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg := config.LoadConfig()
	rdb := redisclient.RedisClient(cfg)
	app := gin.Default()

	conn, err := grpc.NewClient(cfg.AUTH_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authClient := auth.NewAuthServiceClient(conn)

	api := app.Group("/api/v1")
	api.Use(middlewares.RateLimit(rdb))

	public := api.Group("/")
	private := api.Group("/")
	private.Use(middlewares.ValidiateJWT(cfg))

	handlers.RegisterAuthRoutes(public, private, authClient)

	if err := app.Run(":" + cfg.APP_PORT); err != nil {
		log.Println(err)
	}
}
