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
	"os"
	"encoding/csv"
	"path/filepath"
	"runtime"
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

func ReadWordsFile(c *gin.Context, ctx context.Context, client *firestore.Client) error{
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
		var  path = "static/" +  c.Query("set_id") + ".csv"
		file, err := os.Create(path)
		checkError("Cannot create file", err)
		defer file.Close()

		writer := csv.NewWriter(file)


		for _, value := range Words {
			var data = []string{}
			data = append(data,value.Content,value.Language, value.Translation, value.Pronunciation )
			err := writer.Write(data)
			checkError("Cannot write to file", err)
		}
		writer.Flush()
		//c.Redirect(http.StatusMovedPermanently, c.Request.Host + c.Request.URL.Path + path)

		//Words
		//c.JSON(http.StatusCreated, gin.H{
		//	"status": http.StatusCreated,
		//	"message": "Read words successfully!",
		//	"resources": Words,
		//})
		//targetPath := filepath.Abs(path)
		env := os.Getenv("APP_ENV")
		var _filepath string
		_, filename, _, _ := runtime.Caller(1)
		if env == "production" {
			_filepath = filepath.Join(filepath.Dir(filename), "/src/camera-dictionary-golang/static/"+c.Query("set_id") + ".csv")
			log.Println("Running api server in dev mode")
		} else {
			_filepath = filepath.Join(filepath.Dir(filename), "/static/"+c.Query("set_id") + ".csv")
			log.Println("Running api server in production mode")
		}

		log.Print(_filepath)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+ c.Query("set_id") + ".csv" )
		c.Header("Content-Type", "application/octet-stream")
		defer c.File(_filepath)
		///Users/thanhdatvo/Desktop/gogo/src/camera-dictionary-golang/static/result.csv
		///Users/thanhdatvo/Desktop/gogo/src/camera-dictionary-golang/static/
		//c.File("/Users/thanhdatvo/Desktop/gogo/src/camera-dictionary-golang/static/result.csv")
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


func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}