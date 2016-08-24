package main

import (
	"github.com/gohttp/app"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

var (
	MongoAddr = "localhost"
	MongoDB   = "test"
)

func main() {
	os.Environ()
	if addr := os.Getenv("MONGO_ADDR"); len(addr) != 0 {
		MongoAddr = addr
	}
	if db := os.Getenv("MONGO_DB"); len(db) != 0 {
		MongoDB = db
	}
	sess, err := mgo.Dial(MongoAddr)
	if err != nil {
		panic(err)
	}
	shopHandler := &ShopHandler{sess}
	productHandler := &ProductHandler{sess}
	m := app.New()

	m.Post("/shops", shopHandler.Post)
	m.Get("/shops", shopHandler.Get)
	m.Put("/shops/myshop", shopHandler.PutMyShop)
	m.Get("/shops/myshop", shopHandler.GetMyShop)
	m.Get("/shops/:id", shopHandler.GetById)

	m.Post("/products", productHandler.Post)
	m.Get("/products", productHandler.Get)
	m.Put("/products/:id", productHandler.PutById)
	m.Get("/products/:id", productHandler.GetById)
	m.Del("/products/:id", productHandler.DelById)

	if err := http.ListenAndServe(":80", m); err != nil {
		panic(err)
	}
}
