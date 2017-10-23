package controllers

import (
	//"fmt"
	//"log"
	"github.com/gin-gonic/gin"
	//"os"
	"net/http"
	//"errors"
	//"time"
	//
	//"google.golang.org/api/iterator"

	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
	//"fmt"
	//"reflect"
	//"go/doc"
	"google.golang.org/api/iterator"
	"fmt"
	//"go/doc"
 	//"encoding/json"
 	"projects/camera-dictionary/modals"
)
func Login(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	users := client.Collection("users")
	iter := users.Where("username", "==", c.PostForm("username")).Where("password", "==", c.PostForm("password")).Documents(ctx)
	var x int = 0
	var user modals.User
	var id string
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Something happens!"})
			return nil
		}
		x = x + 1
		dsnap.DataTo(&user)
		id = dsnap.Ref.ID
		fmt.Println(dsnap.Data())
	}

	if x == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "Not found!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "Login successfully!",
			"resource":  map[string]interface{}{
				"user_id": id,
			},
		})
	}
	return nil
}

func Register(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Something happens!"})
		return nil
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Registration successfully!",
		})
	return nil
}

