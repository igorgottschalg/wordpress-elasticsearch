package main

import (
	"log"
	"net/http"
	"encoding/json"
	"time"
	"bytes"
	"strconv"
	"os"

	"github.com/gorilla/mux"
	"github.com/fatih/color"
)

var ELASTIC_SEARCH_URL = get_elastic_searchUrl()
const POST_INDEX = `{
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

func get_elastic_searchUrl() string {
    var url string
    if (os.Getenv("ELASTIC_SEARCH_URL") != ""){
        url = os.Getenv("ELASTIC_SEARCH_URL")
    }else{
        url = "http://localhost:9200"
    }
    return url
}

func errorHandler(err error) bool{
    if err != nil {
        color.Red(err.Error())
    }
    return err != nil
}

func save_post_handle(response http.ResponseWriter, req *http.Request) {
    var post Post

    err := json.NewDecoder(req.Body).Decode(&post)
    errorHandler(err)

    post_to_save, err := json.Marshal(post)
    errorHandler(err)

    request,_ := http.NewRequest("PUT", ELASTIC_SEARCH_URL + "/posts/_doc/" + strconv.Itoa(post.ID), bytes.NewBuffer(post_to_save))
    request.Header.Set("Content-Type", "application/json")
    client := &http.Client{}

    _, err = client.Do(request)
    errorHandler(err)

    color.Green("Object saved:")
    log.Println(string(post_to_save))
}

func register_elastic_index(){
    color.White("Creating Elastic Search index")
    req,err :=  http.NewRequest("PUT", ELASTIC_SEARCH_URL + "/posts", bytes.NewBuffer([]byte(POST_INDEX)))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err = client.Do(req)
    errorHandler(err)

    if (err==nil) {
        color.Green("Index created.")
    }

}

func main() {
    color.Blue("Inicializing webserver...")
    color.White("Elastic search url: " + ELASTIC_SEARCH_URL)

    router := mux.NewRouter()
    router.HandleFunc("/", save_post_handle).Methods("POST")

    srv := &http.Server{
        Handler:      router,
        Addr:         "0.0.0.0:3000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    register_elastic_index()

    log.Fatal(srv.ListenAndServe())
}