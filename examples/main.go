package main

import (
	"github.com/MrToy/business"
	"github.com/gohttp/app"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

func main() {
	MongoAddr := "localhost"
	if addr := os.Getenv("MONGO_ADDR"); len(addr) != 0 {
		MongoAddr = addr
	}
	MongoDB := "test"
	if db := os.Getenv("MONGO_DB"); len(db) != 0 {
		MongoDB = db
	}
	sess, err := mgo.Dial(MongoAddr)
	if err != nil {
		panic(err)
	}
	shopHandler := &business.ShopHandler{sess, MongoDB}
	productHandler := &business.ProductHandler{sess, MongoDB}
	orderHandler := &business.OrderHandler{sess, MongoDB}
	payHandler := &business.PayHandler{sess, MongoDB}
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

	m.Put("/orders/confirm/:id", orderHandler.ConfirmById)
	m.Put("/orders/cancle/:id", orderHandler.CancleById)
	m.Put("/orders/deliver/:id", orderHandler.DeliverById)
	m.Post("/orders", orderHandler.Post)
	m.Get("/orders", orderHandler.Get)
	m.Get("/orders/:id", orderHandler.GetById)

	m.Post("/pay", payHandler.Notify)
	m.Get("/pay/:id", payHandler.RedirectById)
	// m.Put("/orders/:id", productHandler.PutById)

	// m.Del("/products/:id", productHandler.DelById)

	if err := http.ListenAndServe(":80", m); err != nil {
		panic(err)
	}
}
