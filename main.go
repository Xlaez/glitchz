package main

import (
	"glitchz/pkg/server"
	"glitchz/pkg/utils"
	"log"

	_ "glitchz/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Docs For Glitchz
// @version         1.0
// @description     This is  the swagger documentation for Glitchz.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5500
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot not load config", err)
	}

	s := server.Run()
	s.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(s.Run(":" + config.Port))
}
