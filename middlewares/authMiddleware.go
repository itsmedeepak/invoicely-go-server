package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"tmp-invoicely.co/utils"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		log.Printf("Method: %s, Endpoint: %s", c.Request.Method, c.Request.URL.Path)

		path := c.Request.URL.Path
		log.Println("mddlle")
		// Skip authentication for auth routes

		log.Println(path)

		if strings.HasPrefix(path, "/auth/") {
			c.Next()
			return
		}

		log.Println(path)

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			c.Abort()
			return
		}

		// Get secret token from environment
		secretToken := utils.GetEnv("SECRET_KEY")

		if secretToken == "" {
			log.Println("SECRET_KEY is missing in environment variables")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretToken), nil
		})

		if err != nil || !token.Valid {
			log.Println("Token verification failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set claims into context
		if userId, exists := claims["_id"]; exists {
			c.Set("user_id", userId)
			log.Println("Authenticated User ID:", userId)
		}
		if email, exists := claims["email"]; exists {
			c.Set("email", email)
		}

		c.Next()
	}
}
