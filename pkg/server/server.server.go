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
	redis_client            *redis.Client
	auth_controller         controllers.AuthController
	user_controller         controllers.UserController
	post_controller         controllers.PostController
	comment_controller      controllers.CommentController
	group_controller        controllers.GroupController
	group_msgs_controller   controllers.GroupMsgsController
	contact_controller      controllers.ContactController
	private_msgs_controller controllers.PrivateMsgsController
)

func InitTokenMaker(config utils.Config) (token.Maker, error) {
	var err error
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create the token maker: %w", err)
	}
	return tokenMaker, nil
}

func initCols(client *mongo.Client, config utils.Config, ctx context.Context, tokenMaker token.Maker, redis_client *redis.Client) (*controllers.AuthController, *controllers.UserController, *controllers.PostController, *controllers.CommentController, *controllers.GroupController, *controllers.ContactController, *controllers.PrivateMsgsController, *controllers.GroupMsgsController) {
	users_col := client.Database(config.DbName).Collection(config.UsersCol)
	tokens_col := client.Database(config.DbName).Collection(config.TokensCol)
	posts_col := client.Database(config.DbName).Collection(config.PostsCol)
	comments_col := client.Database(config.DbName).Collection(config.CommentsCol)
	groups_col := client.Database(config.DbName).Collection(config.GroupsCol)
	group_requests_col := client.Database(config.DbName).Collection(config.GroupRequestCol)
	contacts_col := client.Database(config.DbName).Collection(config.ContactsCol)
	messages_col := client.Database(config.DbName).Collection(config.MsgCol)
	group_msgs_col := client.Database(config.DbName).Collection(config.GroupMsgCol)

	auth_service := services.NewAuthService(users_col, ctx)
	user_service := services.NewUserService(users_col, ctx)
	token_service := services.NewTokenService(tokens_col, tokenMaker, ctx)
	post_service := services.NewPostService(posts_col, ctx)
	comment_service := services.NewCommentService(comments_col, ctx)
	group_service := services.NewGroupService(groups_col, ctx)
	group_requests_service := services.NewGroupRequestService(group_requests_col, ctx)
	contact_service := services.NewContactService(contacts_col, ctx)
	messages_service := services.NewPrivateMsgs(messages_col, ctx)
	group_msgs_service := services.NewGroupMsgs(group_msgs_col, ctx)

	auth_controller = controllers.NewAuthController(auth_service, token_service, tokenMaker, config, redis_client)
	user_controller = controllers.NewUserController(user_service, token_service, tokenMaker, config, redis_client)
	post_controller = controllers.NewPostController(post_service, tokenMaker, config, redis_client)
	comment_controller = controllers.NewCommentController(comment_service, tokenMaker, config, redis_client)
	group_controller = controllers.NewGroupController(group_service, group_requests_service, tokenMaker, config, redis_client)
	contact_controller = controllers.NewContactController(contact_service, user_service, tokenMaker, config, redis_client)
	group_msgs_controller = controllers.NewGroupMsgsController(group_msgs_service, tokenMaker, config, redis_client)
	private_msgs_controller = controllers.NewPrivateMsgsController(messages_service, contact_service, tokenMaker, config, redis_client)

	return &auth_controller, &user_controller, &post_controller, &comment_controller, &group_controller, &contact_controller, &private_msgs_controller, &group_msgs_controller
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

	auth_col, users_col, post_col, comment_col, groups_col, contacts_col, private_msgs_col, group_msg_col := initCols(mongoClient, config, ctx, tokenMaker, redis_client)

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
	routes.PostRoutes(server, *post_col, tokenMaker)
	routes.CommentRoute(server, *comment_col, tokenMaker)
	routes.GroupRoutes(server, *groups_col, *group_msg_col, tokenMaker)
	routes.ContactRoutes(server, *contacts_col, tokenMaker)
	routes.PrivateMsgtRoutes(server, *private_msgs_col, tokenMaker)
	return server
}
