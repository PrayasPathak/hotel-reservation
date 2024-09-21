package api

import (
	"errors"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
	}
	return c.JSON(hotel)
}
