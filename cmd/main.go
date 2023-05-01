package main

import (
	"glitchz/pkg/server"
	"glitchz/pkg/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot not load config", err)
	}

	s := server.Run()
	// s.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(s.Run(":" + config.Port))
}
