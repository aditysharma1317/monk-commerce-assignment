package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"monk-commerce-assignment/config"
	"monk-commerce-assignment/constants"
	"monk-commerce-assignment/handlers"
	"monk-commerce-assignment/utils/db"
	ulog "monk-commerce-assignment/utils/log"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	file, err := os.Open(env + ".json")
	if err != nil {
		log.Println("Unable to open file. Err:", err)
		os.Exit(1)
	}

	var cnf *config.Config
	config.ParseJSON(file, &cnf)
	config.Set(cnf)

	constants.Logger = ulog.New("Global", cnf.AppName, cnf.LogLevel)
	db.Init(&db.Config{
		URL: cnf.DatabaseURL,
	})

	router := gin.Default()
	handlers.SetupRoutes(router)
	router.Run(":8080")
}
