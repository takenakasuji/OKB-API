package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// DBの定義
type Messages struct {
	ID   int64
	Body string
}

// 入力の定義
type okbInput struct {
	Call string
}

// 出力の定義
type okbOutput struct {
	Body string
}

var db gorm.DB

func init() {
	db, _ = gorm.Open("sqlite3", "./okb.db")
}

func okbMessage(w rest.ResponseWriter, req *rest.Request) {
	input := okbInput{}
	err := req.DecodeJsonPayload(&input)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if input.Call != "Hi, OKB" {
		rest.Error(w, "Call format is error", 400)
		return
	}

	// Query用の乱数生成
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(3)

	message := Messages{}
	db.Find(&message, r)

	w.WriteJson(&message)
}

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/api/v1", okbMessage),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
