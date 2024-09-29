package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	*db.Store
	client *mongo.Client
}

func setUp(t *testing.T) *testDB {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	testDBURI := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testDB{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	dbname := os.Getenv(db.MongoDBEnvName)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
