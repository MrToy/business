package main

import (
	//"fmt"
	"encoding/json"
	"github.com/gohttp/app"
	"github.com/gohttp/response"
	"io/ioutil"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"time"
)

type ProductProfile struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	Images    []string      `json:"images"`
	Price     int           `json:"price"`
	Inventory int           `json:"inventory"`
	Active    bool          `json:"active"`
	Date      time.Time     `json:"date"`
}

type ProductEditable struct {
	Name      string   `json:"name"`
	Images    []string `json:"images"`
	Price     int      `json:"price"`
	Info      string   `json:"info"`
	Inventory int      `json:"inventory"`
	Active    bool     `json:"active"`
	Props     []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"props"`
	Spec []struct {
		Name  string   `json:"name"`
		Value []string `json:"value"`
	} `json:"spec"`
}

type Product struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Shop      bson.ObjectId `json:"shop"`
	Name      string        `json:"name"`
	Images    []string      `json:"images"`
	Price     int           `json:"price"`
	Info      string        `json:"info"`
	Inventory int           `json:"inventory"`
	Active    bool          `json:"active"`
	Date      time.Time     `json:"date"`
	Props     []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"props"`
	Spec []struct {
		Name  string   `json:"name"`
		Value []string `json:"value"`
	} `json:"spec"`
}

type ProductMeta struct {
	Total int `json:"total"`
}

func ProductRoute(m *app.App) {
	m.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		var product Product
		if err := json.Unmarshal(data, &product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		product.Id = bson.NewObjectId()
		product.Date = time.Now()
		sess := Sess.Clone()
		defer sess.Close()
		shop := Shop{}
		if err := sess.DB(MongoDB).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		product.Shop = shop.Id
		if err := sess.DB(MongoDB).C("products").Insert(product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		response.NoContent(w)
	})
	m.Put("/products/:id", func(w http.ResponseWriter, r *http.Request) {
		owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		var product ProductEditable
		if err := json.Unmarshal(data, &product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		sess := Sess.Clone()
		defer sess.Close()
		shop := Shop{}
		if err := sess.DB(MongoDB).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if err := sess.DB(MongoDB).C("products").Update(&bson.M{"_id": bson.ObjectIdHex(r.FormValue(":id")), "shop": shop.Id}, &bson.M{"$set": product}); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		response.NoContent(w)
	})
	m.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		skip, _ := strconv.Atoi(r.FormValue("skip"))
		limit, _ := strconv.Atoi(r.FormValue("limit"))
		active, _ := strconv.ParseBool(r.FormValue("active"))
		findQuery := bson.M{}
		if word := r.FormValue("word"); len(word) != 0 {
			findQuery["name"] = &bson.M{"$regex": word, "$options": "i"}
		}
		if shop := r.FormValue("shop"); len(shop) != 0 {
			if !bson.IsObjectIdHex(shop) {
				response.Forbidden(w, "error shop")
				return
			}
			findQuery["shop"] = bson.ObjectIdHex(shop)
			active = false
		}
		sess := Sess.Clone()
		defer sess.Close()
		if active {
			shops := []Shop{}
			if err := sess.DB(MongoDB).C("shops").Find(&bson.M{"active": true}).All(&shops); err != nil {
				response.InternalServerError(w, err.Error())
				return
			}
			var shopIds []bson.ObjectId
			for _, shop := range shops {
				shopIds = append(shopIds, shop.Id)
			}
			findQuery["shop"] = &bson.M{"$in": shopIds}
			findQuery["active"] = true
		}
		products := []ProductProfile{}
		if err := sess.DB(MongoDB).C("products").Find(&findQuery).Skip(skip).Limit(limit).Sort("-date").All(&products); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if ok, _ := strconv.ParseBool(r.FormValue("meta")); ok {
			total, _ := sess.DB(MongoDB).C("products").Find(&findQuery).Count()
			response.JSON(w, &ProductMeta{Total: total})
			return
		}
		response.JSON(w, &products)
	})
	m.Get("/products/:id", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue(":id")
		if !bson.IsObjectIdHex(id) {
			response.Forbidden(w, "error id")
			return
		}
		sess := Sess.Clone()
		defer sess.Close()
		product := Product{}
		if err := sess.DB(MongoDB).C("products").FindId(bson.ObjectIdHex(id)).One(&product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		response.JSON(w, &product)
	})
	m.Del("/products/:id", func(w http.ResponseWriter, r *http.Request) {
		id, owner := r.FormValue(":id"), bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
		if !bson.IsObjectIdHex(id) {
			response.Forbidden(w, "error id")
			return
		}
		sess := Sess.Clone()
		defer sess.Close()
		shop := Shop{}
		if err := sess.DB(MongoDB).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if err := sess.DB(MongoDB).C("products").Remove(bson.M{"_id": bson.ObjectIdHex(id), "shop": shop.Id}); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		response.NoContent(w)
	})
}
