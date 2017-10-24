package controllers

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
	"camera-dictionary-golang/modals"
)

func CreateWord(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	doc, _, err := client.Collection("words").Add(ctx, map[string]interface{}{
		"set_id": c.Query("set_id"),
		"type": c.PostForm("type"),
		"content": c.PostForm("content"),
		"language": c.PostForm("language"),
		"translation": c.PostForm("translation"),
		"pronunciation": c.PostForm("pronunciation"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Something happens!"})
		return nil
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Word item created successfully!",
		"resourceId": doc.ID,
	})
	return nil
}

func ReadWords(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	wordsRef := client.Collection("words")
	iter := wordsRef.Where("set_id", "==", c.Query("set_id")).Documents(ctx)
	var x int = 0
	var Words = []modals.Word{}
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
		var s modals.Word
		dsnap.DataTo(&s)
		s.Id = dsnap.Ref.ID
		Words = append(Words, s)
	}

	if x == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "Not found!",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status": http.StatusCreated,
			"message": "Read words successfully!",
			"resources": Words,
		})
	}
	return nil
}

func ReadWord(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	wordsRef := client.Collection("words")
	dsnap, err := wordsRef.Doc(c.Param("wordId")).Get(ctx)
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
	var word modals.Word
	dsnap.DataTo(&word)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Read word successfully!",
		"resources": word,
	})
	return nil
}

func UpdateWord(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	wordsRef := client.Collection("words").Doc(c.Param("wordId"))
	batch := client.Batch()
	batch.Set(wordsRef, map[string]interface{}{
		"set_id": c.PostForm("set_id"),
		"type": c.PostForm("type"),
		"content": c.PostForm("content"),
		"language": c.PostForm("language"),
		"translation": c.PostForm("translation"),
		"pronunciation": c.PostForm("pronunciation"),
	})
	_, err := batch.Commit(ctx)
	//dsnap, err := wordsRef..Get(ctx)
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
	//var word Word
	//dsnap.DataTo(&word)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Update word successfully!",
	})
	return nil
}

func DeleteWord(c *gin.Context, ctx context.Context, client *firestore.Client) error{
	_, err := client.Collection("words").Doc(c.Param("wordId")).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Delete word successfully!",
	})
	return nil
}



