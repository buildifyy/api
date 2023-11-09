package main

import (
	"api/pkg/common"
	"api/pkg/db"
	"api/pkg/instance"
	"api/pkg/template"
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	dbConnectionString := os.Getenv("ConnectionString")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbConnectionString).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println("error connecting database: ", err)
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println("error disconnecting database: ", err)
			panic(err)
		}
	}()

	dbRepository := db.NewRepository(client)

	if err := dbRepository.Ping(); err != nil {
		log.Println("error pinging database: ", err)
		panic(err)
	}
	log.Println("successfully connected to database")

	r := gin.Default()

	r.Use(cors.Default())

	commonService := common.NewService(dbRepository)
	commonController := common.NewController(commonService)
	common.RegisterRoutes(r, commonController)

	templateService := template.NewService(dbRepository)
	templateController := template.NewController(templateService)
	template.RegisterRoutes(r, templateController)

	instanceService := instance.NewService(dbRepository, templateService, commonService)
	instanceController := instance.NewController(instanceService)
	instance.RegisterRoutes(r, instanceController)

	if err = r.Run(); err != nil {
		panic(err)
	}
}
