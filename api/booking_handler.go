package api

import (
	"errors"
	"net/http"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// This needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

// This needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}

	if booking.ID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type:    "error",
			Message: "not authorized",
		})
	}
	return c.JSON(booking)
}
