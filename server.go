package main

import (
	"log"
	"net/http"
	"encoding/json"
	"time"
	"bytes"
	"strconv"

	"github.com/gorilla/mux"
)

var elasticSearchUrl = "http://elasticsearch:9200"
var postIndex = `{
  "index": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "_id": {
        "type": "interger"
      },
      "fields": {
        "properties": {
          "id": {
            "type": "interger",
            "index": true
          },
          "content": {
            "type": "text",
            "index": true
          },
          "name": {
            "type": "text",
            "index": true
          },
          "image": {
            "type": "text"
          },
          "url": {
            "type": "text"
          },
          "posttype": {
            "type": "text"
          },
          "keywords": {
            "type": "nested",
            "index": true
          }
        }
      }
    }
  }
}`

type Post struct {
    ID         int    `json:"id"`
    Name       string `json:"name"`
    Content    string `json:"content"`
    Image      string `json:"image"`
    Url        string `json:"url"`
    PostType   string `json:"post_type"`
    Keywords []string `json:"keywords"`
}

func savePostHandleFunc (response http.ResponseWriter, req *http.Request) {
    response.Header().Add("Content-Type", "application/json")
    var post Post

    err := json.NewDecoder(req.Body).Decode(&post)
    if err != nil {
        log.Println(err)
        response.Header().Set("Connection", "close")
        return
    }

    post_to_save, err := json.Marshal(post)
    if err != nil {
        log.Println(err)
        response.Header().Set("Connection", "close")
        return
    }

    request,_ := http.NewRequest("PUT", elasticSearchUrl+"/posts/_doc/"+strconv.Itoa(post.ID), bytes.NewBuffer(post_to_save))
    request.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    _, err = client.Do(request)
    if err != nil {
        log.Println(err)
        response.Header().Set("Connection", "close")
        return
    }

    log.Println(string(post_to_save))
}

func registerElasticIndex(){
    req,err :=  http.NewRequest("PUT", elasticSearchUrl+"/posts", bytes.NewBuffer([]byte(postIndex)))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err = client.Do(req)
    if err != nil {
        log.Println(err)
    }
}

func main() {
    registerElasticIndex()

    router := mux.NewRouter()
    router.HandleFunc("/", savePostHandleFunc).Methods("POST")

    srv := &http.Server{
        Handler:      router,
        Addr:         "0.0.0.0:3030",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}