package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


var jwtKey = []byte("rahasia-mei")

func GenerateToken(username string)(string, error){
	claims := jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString(jwtKey)
}
func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == ""{
			c.JSON(401,gin.H{"error": "Token tidak ada"})
			c.Abort()
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr,"Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token)(interface{},error){
			return jwtKey, nil
		})
		if err != nil || !token.Valid{
			c.JSON(401,gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}
		c.Set("username",claims["username"])
		c.Next()
	}
}