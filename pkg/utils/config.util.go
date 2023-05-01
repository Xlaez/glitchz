package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	CloudinaryKey       string        `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryName      string        `mapstructure:"CLOUDINARY_NAME"`
	CloudinarySecret    string        `mapstructure:"CLOUDINARY_SECRET_KEY"`
	TokenKey            string        `mapstructure:"TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MongoUri            string        `mapstructure:"MONGODB_URI"`
	Port                string        `mapstructure:"port"`
	DbName              string        `mapstructure:"DB_NAME"`
	RedisUri            string        `mapstructure:"REDIS_URL"`
	SmtpHost            string        `mapstructure:"SMTP_HOST"`
	SmtpPort            string        `mapstructure:"SMTP_PORT"`
	SmtpUsername        string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword        string        `mapstructure:"SMTP_PASSWORD"`
	AppEmail            string        `mapstructure:"EMAIL_FROM"`
	UsersCol            string        `mapstructure:"USERS_COl"`
	ContactsCol         string        `mapstructure:"CONTACTS_COL"`
	PostsCol            string        `mapstructure:"POSTS_COL"`
	CommentsCol         string        `mapstructure:"COMMENTS_COL"`
	MsgCol              string        `mapstructure:"MSG_COL"`
	GroupsCol           string        `mapstructure:"GROUPS_COL"`
	GroupMsgCol         string        `mapstrucutre:"GROUP_MSG_COL"`
	GroupRequestCol     string        `mapstructure:"GROUP_REQUESTS_COL"`
	TokensCol           string        `mapstructure:"TOKEN_COL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
