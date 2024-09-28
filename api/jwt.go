package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PrayasPathak/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if token == "" {
			fmt.Println("Token not present in the header")
			return ErrUnauthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expires := claims["expires"].(float64)
		if time.Now().Unix() > int64(expires) {
			return NewError(http.StatusUnauthorized, ErrTokenExpiredMsg)
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}
		// Set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnauthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse token JWT TOKEN ", err)
		return nil, ErrUnauthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token ", err)
		return nil, ErrUnauthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthorized()
	}
	return claims, nil
}
