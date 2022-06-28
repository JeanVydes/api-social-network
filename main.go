package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	App        *gin.Engine
	Production = false
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	envFilename := "development.env"
	if Production {
		envFilename = "production.env"
	}

	err := godotenv.Load(envFilename)
	if err != nil {
		panic(err)
	}

	InitializeMongoConnection()
	Initialize()
}

func Initialize() {
	App = gin.Default()

	Configure()
	Listen()
}

func Configure() {
	SetRouters()
}

func Listen() {
	App.Run(":" + os.Getenv("PORT"))
}
