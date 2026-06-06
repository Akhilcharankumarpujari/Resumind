package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)


type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	GoogleID string `json:"google_id"`
	jwt.RegisteredClaims
}


type GoogleLoginRequest struct {
	IDToken string `json:"idToken"`
}


func HandleGoogleAuth(c *fiber.Ctx) error {
	var body GoogleLoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Google Client ID not configured on server"})
	}


	payload, err := idtoken.Validate(context.Background(), body.IDToken, clientID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Google token"})
	}


	googleID := payload.Subject
	email := fmt.Sprintf("%v", payload.Claims["email"])
	name := fmt.Sprintf("%v", payload.Claims["name"])
	picture := fmt.Sprintf("%v", payload.Claims["picture"])


	var user User
	result := DB.Where("google_id = ?", googleID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		user = User{
			GoogleID: googleID,
			Email:    email,
			Name:     name,
			Picture:  picture,
		}
		if err := DB.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
		}
	} else if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	} else {

		DB.Model(&user).Updates(map[string]interface{}{
			"name":    name,
			"picture": picture,
		})
	}


	jwtToken, err := generateJWT(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate session token"})
	}


	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
		},
	})
}


func HandleGetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*JWTClaims)
	return c.JSON(fiber.Map{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
	})
}


func HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}


func AuthMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Cookies("auth_token")
	if tokenStr == "" {

		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Not authenticated"})
	}

	claims, err := validateJWT(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired session"})
	}

	c.Locals("user", claims)
	return c.Next()
}


func generateJWT(user User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-this-to-a-secure-random-secret-in-production"
	}

	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Picture:  user.Picture,
		GoogleID: user.GoogleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}


func validateJWT(tokenStr string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-this-to-a-secure-random-secret-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
