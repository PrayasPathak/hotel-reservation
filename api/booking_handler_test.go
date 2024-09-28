package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PrayasPathak/hotel-reservation/db/fixtures"
	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "ramesh", "sharma", false)
		hotel     = fixtures.AddHotel(tdb.Store, "bar hotel", "bermuda", 5, nil)
		rooms     = fixtures.AddRoom(tdb.Store, "large", true, 299.99, hotel.ID)
		from      = time.Now()
		till      = from.AddDate(0, 0, 7)
		booking   = fixtures.AddBooking(tdb.Store, rooms.ID, user.ID, from, till)
		app       = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		admin          = app.Group("/", JWTAuthentication(tdb.User), AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 response but got %d", resp.StatusCode)
	}
	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got none")
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s but got %s", booking.UserID, have.UserID)
	}

	// Non-admin user cannot access
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(tdb.Store, "jimmy", "smith", false)
		user           = fixtures.AddUser(tdb.Store, "ramesh", "sharma", false)
		hotel          = fixtures.AddHotel(tdb.Store, "bar hotel", "bermuda", 5, nil)
		rooms          = fixtures.AddRoom(tdb.Store, "large", true, 299.99, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 7)
		booking        = fixtures.AddBooking(tdb.Store, rooms.ID, user.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected non 200 code but got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, bookingResp.ID)
	}

	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s but got %s", booking.UserID, bookingResp.UserID)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 code but got %d", resp.StatusCode)
	}
}
