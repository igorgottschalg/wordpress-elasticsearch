package main

import (
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

type POST struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	URL      string `json:"url"`
	PostType string `json:"post_type"`
}

func savePostHandleFunc (w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "application/json")
    var post POST

    err := json.NewDecoder(req.Body).Decode(&post)
    if err != nil {
        log.Println(err)
        return
    }

    post_to_save, err := json.Marshal(&post)
    if err != nil {
        return
    }

    client := redisClient()
    err = client.Set(string(post.ID), string(post_to_save), 0).Err()
    if err != nil {
        log.Println(err)
    }
}

func redisClient() *redis.Client{
    client := redis.NewClient(&redis.Options{
   		Addr:   "127.0.0.1:6379",
   		Password: "",
   		DB: 0,
   	})
   	return client
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", savePostHandleFunc).Methods("POST")
    srv := &http.Server{
        Handler:      router,
        Addr:         "127.0.0.1:3030",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}