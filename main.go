package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)


var oauthConf *oauth2.Config

func init(){
	_ = godotenv.Load()

	oauthConf = &oauth2.Config{
		ClientID: os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL: "http://localhost:8080/auth/github/callback",
		Scopes: []string{"user:email"},
		Endpoint: github.Endpoint,
	}
}
func main() {
	r := gin.Default()
	r.GET("/auth/github",func(c *gin.Context) {
		url := oauthConf.AuthCodeURL("randomstate")
		c.Redirect(http.StatusTemporaryRedirect,url)
	})
	r.GET("/auth/github/callback",func(c *gin.Context) {
		code := c.Query("code")
		if code == ""{
			c.JSON(400,gin.H{"error": "No code"})
			return
		}
		token, err := oauthConf.Exchange(context.Background(),code)
		if err != nil{
			c.JSON(500,gin.H{"error": "Token exchange failed"})
			return
		}
		client := oauthConf.Client(context.Background(),token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil{
			c.JSON(500,gin.H{"error": "Failed to fetch user"})
			return
		}
		defer resp.Body.Close()

		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil{
			c.JSON(500,gin.H{"error": "Failed to decode user"})
			return
		}

		// Ambil username GitHub
		username := fmt.Sprintf("%v", userInfo["login"])

		// Generate JWT
		jwtToken, err := GenerateToken(username)
		if err != nil{
			c.JSON(500, gin.H{"error": "Failed to generate JWT"})
			return
		}

		c.SetCookie("auth_token", jwtToken, 3600 * 24, "/", "localhost", false, true)
		c.JSON(200, gin.H{"message": "Login via GitHub berhasil", "username": username})

	})
	r.GET("/profile",AuthMiddleware(),func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(200,gin.H{
			"message": "hello " + username,
		})
	})
	r.Run(":8080")
}