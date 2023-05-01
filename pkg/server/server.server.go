package server

import (
	"context"
	"fmt"
	"glitchz/pkg/controllers"
	"glitchz/pkg/routes"
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	cors "github.com/rs/cors/wrapper/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	redis_client    *redis.Client
	auth_controller controllers.AuthController
	user_controller controllers.UserController
)

func InitTokenMaker(config utils.Config) (token.Maker, error) {
	var err error
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create the token maker: %w", err)
	}
	return tokenMaker, nil
}

func initCols(client *mongo.Client, config utils.Config, ctx context.Context, tokenMaker token.Maker, redis_client *redis.Client) (*controllers.AuthController, *controllers.UserController) {
	users_col := client.Database(config.DbName).Collection(config.UsersCol)
	tokens_col := client.Database(config.DbName).Collection(config.TokensCol)

	auth_service := services.NewAuthService(users_col, ctx)
	user_service := services.NewUserService(users_col, ctx)

	auth_controller = controllers.NewAuthController(auth_service, tokenMaker, config, *tokens_col, redis_client)
	user_controller = controllers.NewUserController(user_service, tokenMaker, config, *tokens_col, redis_client)

	return &auth_controller, &user_controller
}

func Run() *gin.Engine {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load env", err)
	}

	ctx := context.TODO()
	tokenMaker, err := InitTokenMaker(config)

	if err != nil {
		log.Panic((err.Error()))
	}

	mongoConn := options.Client().ApplyURI(config.MongoUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		log.Panic((err.Error()))
	}

	redis_client = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redis_client.Ping(ctx).Result(); err != nil {
		log.Panic(err.Error())
	}

	err = redis_client.Set(ctx, "test", "Redis on!", 0).Err()

	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println("Redis client connected successfully!")

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic((err.Error()))
	}

	fmt.Println("MongoDB connection succesful!")

	auth_col, users_col := initCols(mongoClient, config, ctx, tokenMaker, redis_client)

	server := gin.Default()
	server.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		Debug:            true, // remeber to off this for prod
		AllowedMethods:   []string{"POST", "GET", "PATCH", "DELETE", "PURGE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           3,
	}))

	// defer mongoClient.Disconnect(ctx)

	routes.AuthRoutes(server, *auth_col, tokenMaker)
	routes.UserRoutes(server, *users_col, tokenMaker)
	// routes.PoductRoutes(server, *prod_col, tokenMaker)

	return server
}
