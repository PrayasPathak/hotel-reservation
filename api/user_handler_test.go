package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	UserStore db.UserStore
}

func setUp(t *testing.T) *testDB {
	const testDBURI = "mongodb://localhost:27017/hotel_reservation_test"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBURI))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to test database:", testDBURI)
	return &testDB{
		UserStore: db.NewMongoUserStore(client),
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	testDB := setUp(t)
	defer testDB.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testDB.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		FirstName: "Elka",
		LastName:  "Doe",
		Email:     "elka@example.com",
		Password:  "smithelka",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Error("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expecting the encrypted password not to be included in json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s, but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s, but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s, but got %s", params.Email, user.Email)
	}
}
