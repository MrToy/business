package main

import (
	"encoding/json"
	"github.com/gohttp/app"
	"github.com/gohttp/response"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Order struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Shop    bson.ObjectId `json:"shop"`
	Buyer   bson.ObjectId `json:"buyer"`
	Contact Contact       `json:"contact"`
	State   string        `json:"state"`
	Date    time.Time     `json:"date"`
	Price   int           `json:"price"`
	Items   []struct {
		Id     bson.ObjectId     `json:"id"`
		Spec   map[string]string `json:"spec"`
		Number int               `json:"number"`
	} `json:"products"`
}

func OrderRoute(m *app.App) {
	m.Post("/order", func(w http.ResponseWriter, r *http.Request) {
		buyer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		var order Order
		if err := json.Unmarshal(data, &order); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		order.Id = bson.NewObjectId()
		order.Date = time.Now()
		order.Buyer = buyer
		order.State = "待付款"
		order.Price = 0
		sess := Sess.Clone()
		defer sess.Close()
		for _, item := range order.Items {
			product := Product{}
			if err := sess.DB(MongoDB).C("products").FindId(item.Id).One(&product); err != nil {
				response.InternalServerError(w, err.Error())
				return
			}
			if product.Shop != order.Shop {
				response.Forbidden(w, "店铺信息错误")
				return
			}
			if product.Inventory < item.Number {
				response.Forbidden(w, "库存不足")
				return
			}
			order.Price += product.Price * item.Number
		}
		if err := sess.DB(MongoDB).C("orders").Insert(order); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		response.NoContent(w)
	})
}
