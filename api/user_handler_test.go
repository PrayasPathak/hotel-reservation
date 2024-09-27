package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	testDB := setUp(t)
	defer testDB.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testDB.User)
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
