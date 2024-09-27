package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/PrayasPathak/hotel-reservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.Store, "ramesh", "sharma", false)
	insertedUser.EncryptedPassword = ""
	// Set the encrypted password to "", because we do NOT return that in any
	// JSON response
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)
	authParams := AuthParams{
		Email:    "ramesh@sharma.com",
		Password: "ramesh_sharma",
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "ramesh", "sharma", false)
	// Set the encrypted password to "", because we do NOT return that in any
	// JSON response
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)
	authParams := AuthParams{
		Email:    "ramesh@sharma.com",
		Password: "sramesh12",
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genericResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("expected generic response type be error, but got %s", genericResp.Type)
	}
	if genericResp.Message != "invalid credentials" {
		t.Fatalf("expected generic response message to be <invalid credentials>, but got %s", genericResp.Type)
	}
}
