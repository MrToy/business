package business

import (
	"github.com/gohttp/response"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
	"net/http"
	"time"
)

type Product struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty" schema:"-"`
	Shop      bson.ObjectId `json:"shop" bson:",omitempty" schema:"-"`
	Date      time.Time     `json:"date" bson:",omitempty" schema:"-"`
	Name      string        `json:"name" validate:"min=1,max=40"`
	Images    []string      `json:"images"`
	Price     int           `json:"price"`
	Inventory int           `json:"inventory"`
	Active    bool          `json:"active"`
	Info      string        `json:"info,omitempty"`
	Props     []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"props,omitempty"`
	Spec []struct {
		Name  string   `json:"name"`
		Value []string `json:"value"`
	} `json:"spec,omitempty"`
}

type ProductQuery struct {
	Skip   int               `bson:"-"`
	Limit  int               `bson:"-"`
	Sort   string            `bson:"-"`
	Word   string            `bson:"-"`
	Shop   string            `bson:"-"`
	Active bool              `bson:",omitempty"`
	ShopId bson.ObjectId     `bson:"shop,omitempty" schema:"-"`
	Name   map[string]string `bson:",omitempty" schema:"-"`
}

type Products struct {
	Total int       `json:"total"`
	Data  []Product `json:"data"`
}

type ProductHandler struct {
	Sess *mgo.Session
	Db   string
}

func (this *ProductHandler) Post(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	r.ParseForm()
	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	product := new(Product)
	schema.NewDecoder().Decode(product, r.Form)
	if err := validator.Validate(product); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	shop := Shop{}
	if err := sess.DB(this.Db).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	product.Shop = shop.Id
	product.Id = bson.NewObjectId()
	product.Date = time.Now()
	if err := sess.DB(this.Db).C("products").Insert(product); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.NoContent(w)
}

func (this *ProductHandler) PutById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	r.ParseForm()
	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	product := new(Product)
	schema.NewDecoder().Decode(product, r.Form)
	if err := validator.Validate(product); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	shop := Shop{}
	if err := sess.DB(this.Db).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if err := sess.DB(this.Db).C("products").Update(&bson.M{"_id": bson.ObjectIdHex(r.FormValue(":id")), "shop": shop.Id}, &bson.M{"$set": product}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.NoContent(w)
}

func (this *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	r.ParseForm()
	query := new(ProductQuery)
	schema.NewDecoder().Decode(query, r.Form)
	if len(query.Sort) == 0 {
		query.Sort = "-date"
	}
	if len(query.Word) != 0 {
		query.Name = map[string]string{"$regex": query.Word, "$options": "i"}
	}
	if len(query.Shop) != 0 && bson.IsObjectIdHex(query.Shop) {
		query.ShopId = bson.ObjectIdHex(query.Shop)
	}
	products := new(Products)
	if err := sess.DB(this.Db).C("products").Find(&query).Select(&bson.M{"info": 0, "props": 0, "spec": 0}).Skip(query.Skip).Limit(query.Limit).Sort(query.Sort).All(&products.Data); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	products.Total, _ = sess.DB(this.Db).C("products").Find(&query).Count()
	response.JSON(w, &products)
}

func (this *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := r.FormValue(":id")
	if !bson.IsObjectIdHex(id) {
		response.Forbidden(w, "error id")
		return
	}
	product := Product{}
	if err := sess.DB(this.Db).C("products").FindId(bson.ObjectIdHex(id)).One(&product); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.JSON(w, &product)
}

func (this *ProductHandler) DelById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id, owner := r.FormValue(":id"), bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	if !bson.IsObjectIdHex(id) {
		response.Forbidden(w, "error id")
		return
	}
	shop := Shop{}
	if err := sess.DB(this.Db).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if err := sess.DB(this.Db).C("products").Remove(bson.M{"_id": bson.ObjectIdHex(id), "shop": shop.Id}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.NoContent(w)
}
