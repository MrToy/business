package main

import (
	"github.com/gohttp/response"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
	"net/http"
	"time"
)

type Shop struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty" schema:"-"`
	Owner  bson.ObjectId `json:"owner" bson:",omitempty" schema:"-"`
	Date   time.Time     `json:"date" bson:",omitempty" schema:"-"`
	Name   string        `json:"name" validate:"min=3,max=40"`
	Active bool          `json:"active"`
	Info   string        `json:"info,omitempty"`
}

type ShopQuery struct {
	Skip   int    `bson:"-"`
	Limit  int    `bson:"-"`
	Sort   string `bson:"-"`
	Active bool   `bson:",omitempty"`
}

type Shops struct {
	Total int    `json:"total"`
	Data  []Shop `json:"data"`
}

type ShopHandler struct {
	sess *mgo.Session
}

func (this *ShopHandler) Post(w http.ResponseWriter, r *http.Request) {
	sess := this.sess.Clone()
	defer sess.Close()
	r.ParseForm()
	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	shop := new(Shop)
	schema.NewDecoder().Decode(shop, r.Form)
	if err := validator.Validate(shop); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	shop.Id = bson.NewObjectId()
	shop.Owner = owner
	shop.Date = time.Now()
	if err := sess.DB(MongoDB).C("shops").EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if err := sess.DB(MongoDB).C("shops").EnsureIndex(mgo.Index{Key: []string{"owner"}, Unique: true}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if err := sess.DB(MongoDB).C("shops").Insert(shop); err != nil {
		if mgo.IsDup(err) {
			response.InternalServerError(w, "店铺已存在")
		} else {
			response.InternalServerError(w, err.Error())
		}
		return
	}
	if err := sess.DB(MongoDB).C("users").Update(&bson.M{"_id": owner}, &bson.M{"$set": &bson.M{"group.saler": true}}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.NoContent(w)
}

func (this *ShopHandler) Get(w http.ResponseWriter, r *http.Request) {
	sess := this.sess.Clone()
	defer sess.Close()
	r.ParseForm()
	shopQuery := new(ShopQuery)
	schema.NewDecoder().Decode(shopQuery, r.Form)
	if len(shopQuery.Sort) == 0 {
		shopQuery.Sort = "-date"
	}
	shops := new(Shops)
	if err := sess.DB(MongoDB).C("shops").Find(&shopQuery).Select(&bson.M{"info": 0}).Sort(shopQuery.Sort).Skip(shopQuery.Skip).Limit(shopQuery.Limit).All(&shops.Data); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	shops.Total, _ = sess.DB(MongoDB).C("shops").Find(&shopQuery).Count()
	response.JSON(w, shops)
}

func (this *ShopHandler) GetMyShop(w http.ResponseWriter, r *http.Request) {
	sess := this.sess.Clone()
	defer sess.Close()
	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	shop := Shop{}
	if err := sess.DB(MongoDB).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.JSON(w, shop)
}

func (this *ShopHandler) PutMyShop(w http.ResponseWriter, r *http.Request) {
	sess := this.sess.Clone()
	defer sess.Close()
	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	r.ParseForm()
	shop := new(Shop)
	schema.NewDecoder().Decode(shop, r.Form)
	if err := validator.Validate(shop); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	if err := sess.DB(MongoDB).C("shops").Update(bson.M{"owner": owner}, bson.M{"$set": shop}); err != nil {
		if mgo.IsDup(err) {
			response.InternalServerError(w, "店铺名已存在")
		} else {
			response.InternalServerError(w, err.Error())
		}
		return
	}
	response.NoContent(w)
}

func (this *ShopHandler) GetById(w http.ResponseWriter, r *http.Request) {
	sess := this.sess.Clone()
	defer sess.Close()
	id := r.FormValue(":id")
	if !bson.IsObjectIdHex(id) {
		response.Forbidden(w, "error id")
		return
	}
	shop := Shop{}
	if err := sess.DB(MongoDB).C("shops").FindId(bson.ObjectIdHex(id)).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.JSON(w, shop)
}
