package middleware

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	// JwtMiddleware       *jwt.HertzJWTMiddleware
	IdentityKey         = "sub"
	PublicKeyAuthServer = "-----BEGIN PUBLIC KEY-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAodxwFdiFKWTG/ZU7vXPdk8ox+nNU\n1JmxsmI8i8tYrYf6QxmwBz13jS/PZsb8dJbMFY3YTMMih6SKz7e+cQ68IbgA7BnY\n5fYFQET4SNHVX/zaH6J70ERJLsRrarmWSXsNbMbnqXlIkoorYXeAn9vsLbr/RPw9\nDYaoq4JrQ+OGsc4LHMw=\n-----END PUBLIC KEY-----\n"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		publicKeyBlock, _ := pem.Decode([]byte(PublicKeyAuthServer))
		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		ECDSAPubKey := publicKey.(*ecdsa.PublicKey)
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return ECDSAPubKey, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(*jwt.MapClaims)
		userID, ok := (*claims)["sub"].(string)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !c.IsAborted() {
			c.Set("userID", userID)
		}
		c.Next()

	}
}
