package api

import (
	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	users := []types.User{
		{FirstName: "John", LastName: "Smith"},
		{FirstName: "John", LastName: "Doe"},
		{FirstName: "Jane", LastName: "Doe"},
		{FirstName: "Mike", LastName: "Daniel"},
	}
	return c.JSON(map[string][]types.User{
		"users": users,
	})
}

func HandleGetUser(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Matt",
		LastName:  "Henry",
	}
	return c.JSON(map[string]types.User{
		"user": user,
	})
}
