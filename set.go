package main

import (
	//"fmt"
	"log"
	"github.com/gin-gonic/gin"
	//"os"
	"net/http"
	//"errors"
	//"time"
	//
	//"google.golang.org/api/iterator"
	//spb "google.golang.org/genproto/googleapis/rpc/status"
	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
	//"fmt"
	//"reflect"
	//"go/doc"
	//"encoding/json"
	"google.golang.org/api/iterator"
	//"fmt"
	//"go/doc"
 	//"encoding/json"
	//"google.golang.org/grpc/status"
	//"golang.org/x/oauth2/clientcredentials"
	//"bytes"}Z
	"google.golang.org/grpc/status"
	//"bytes"
	//"bytes"
	//"strings"
)
func CreateSet(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	doc, _, err := client.Collection("sets").Add(ctx, map[string]interface{}{
		"user_id": c.Query("user_id"),
		"name": c.PostForm("name"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Something happens!"})
		return nil
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Set item created successfully!",
		"resourceId": doc.ID,
	})
	return nil
}

func ReadSets(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	setsRef := client.Collection("sets")
	iter := setsRef.Where("user_id", "==", c.Query("user_id")).Documents(ctx)
	var x int = 0
	var Sets = []Set{}
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
		var s Set
		dsnap.DataTo(&s)
		s.Id = dsnap.Ref.ID
		Sets = append(Sets, s)
	}

	if x == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "Not found!",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status": http.StatusCreated,
			"message": "Read sets successfully!",
			"resources": Sets,
		})
	}
	return nil
}

func ReadSet(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	setsRef := client.Collection("sets")
	dsnap, err := setsRef.Doc(c.Param("setId")).Get(ctx)
	if err != nil {
		log.Printf("Cannot prepare query docs: %#v", err)
		if e, ok := status.FromError(err); ok {
			if e.Code() == 5 {
				//data = bytes.TrimPrefix(data, []byte(""))
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": e.Message()})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
		}
		return nil
	}
	var set Set
	dsnap.DataTo(&set)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Read set successfully!",
		"resources": set,
	})
	return nil
}

func UpdateSet(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	setsRef := client.Collection("sets").Doc(c.Param("setId"))
	batch := client.Batch()
	batch.Set(setsRef, map[string]interface{}{
		"name": c.PostForm("name"),
	})
	_, err := batch.Commit(ctx)
	//dsnap, err := setsRef..Get(ctx)
	if err != nil {
		log.Printf("Cannot prepare query docs: %#v", err)
		if e, ok := status.FromError(err); ok {
			if e.Code() == 5 {
				//data = bytes.TrimPrefix(data, []byte(""))
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": e.Message()})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
		}
		return nil
	}
	//var set Set
	//dsnap.DataTo(&set)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Update set successfully!",
	})
	return nil
}

func DeleteSet(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	_, err := client.Collection("sets").Doc(c.Param("setId")).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Delete set successfully!",
	})
	return nil
}



