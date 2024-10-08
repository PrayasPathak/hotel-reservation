package api

import (
	"errors"
	"net/http"

	"github.com/PrayasPathak/hotel-reservation/db"
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
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

// This needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound("booking")
		}
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{
		"canceled": true,
	}); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(genericResp{
		Type:    "msg",
		Message: "OK",
	})
}
