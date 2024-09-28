package api

import (
	"errors"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var pagination db.Pagination
	if err := c.QueryParser(&pagination); err != nil {
		return ErrBadRequest()
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &pagination)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(pagination.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound("hotel")
		}
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := bson.M{
		"hotelID": oid,
	}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}
