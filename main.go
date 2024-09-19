package main

import (
	"context"
	"flag"
	"log"

	"github.com/PrayasPathak/hotel-reservation/api"
	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri          = "mongodb://localhost:27017"
	dbName         = "hotel-reservation"
	userCollection = "users"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "Port at which the api listens")
	flag.Parse()

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)

	app.Listen(*listenAddr)
}
