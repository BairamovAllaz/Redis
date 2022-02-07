package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)
type Userstruct struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company"`
}
var cach = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})
var ctx = context.Background()

func Getsingel(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	cacherr := cach.Set(ctx, id, body, 10*time.Second).Err()
	if cacherr != nil {
		fmt.Printf("error: %s", cacherr)
		c.Status(http.StatusBadRequest)
		return
	}
	user := Userstruct{}
	marsherr := json.Unmarshal(body, &user)
	if marsherr != nil {
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
func verifyCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		val, err := cach.Get(ctx, id).Bytes()
		if err != nil {
			c.Next()
		}
		user := Userstruct{}
		marsherr := json.Unmarshal(val, &user)
		if marsherr != nil {
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				c.Status(http.StatusBadRequest)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"Cach": user,
		})
	}
}
func main() {
	router := gin.Default()
	router.GET("/:id", verifyCache(), Getsingel)
	router.Run(":8000")
}
