package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/PrayasPathak/hotel-reservation/api"
	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/PrayasPathak/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	var (
		ctx             = context.Background()
		err             error
		mongodbEndpoint = os.Getenv("MONGO_DB_URL")
		dbname          = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbname).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "mike", "daniel", false)
	fmt.Println("Mike -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("Admin -> ", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 299.99, hotel.ID)
	booking := fixtures.AddBooking(store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Printf("Booking -> %+v\n", booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Hotel %d", i)
		location := fmt.Sprintf("Location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
