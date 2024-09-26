package api

import (
	"net/http"

	"github.com/PrayasPathak/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type:    "error",
			Message: "unauthorized",
		})
	}
	return user, nil
}
