package main

import (
	"log"
	"net/http"
	"time"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gorilla/mux"
)

type Fields struct {
    ID         int    `json:"id"`
    Content    string `json:"content"`
    Image      string `json:"image"`
    Url        string `json:"url"`
    PostType   string `json:"post_type"`
    Keywords []string `json:"keywords"`
}

type Post struct {
	ID  int       `json:"id"`
    Fields Fields `json:"fields"`
}

func savePostHandleFunc (w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "application/json")
    var fields Fields

    err := json.NewDecoder(req.Body).Decode(&fields)
    if err != nil {
        log.Println(err)
        return
    }

    post := Post{
        ID: fields.ID,
        Fields: fields,
    }

    post_to_save, err := json.Marshal(post)
    if err != nil {
        log.Println(err)
        return
    }

    ctx := context.Background()

    var (
        docMap map[string]interface{}
    )

    client := esClient()
    client.Set(strconv.Itoa(post.ID), string(post_to_save), 0)
    log.Println(string(post_to_save))
}

func esClient(){
    cfg := elasticsearch.Config{
      Addresses: []string{
        "http://localhost:9200",
      },
    }
    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
      log.Fatalf("Error creating the client: %s", err)
    }
    res, err := es.Info()
    if err != nil {
      log.Fatalf("Error getting response: %s", err)
    }

    log.Println(res)
   	return elasticsearch
}

func main() {
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