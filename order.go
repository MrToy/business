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

type Order struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty" schema:"-"`
	ShopId  bson.ObjectId `json:"shop" bson:",omitempty" schema:"-"`
	Buyer   bson.ObjectId `json:"buyer" bson:",omitempty" schema:"-"`
	State   string        `json:"state" bson:",omitempty" schema:"-"`
	Date    time.Time     `json:"date" schema:"-"`
	Price   int           `json:"price" schema:"-"`
	Shop    string        `json:"-" bson:"-" validate:"len=24"`
	Contact Contact       `json:"contact"`
	Items   []struct {
		Id     bson.ObjectId     `json:"id"`
		Spec   map[string]string `json:"spec"`
		Number int               `json:"number"`
	} `json:"products"`
}

type OrderQuery struct {
	Skip    int           `bson:"-"`
	Limit   int           `bson:"-"`
	Sort    string        `bson:"-"`
	Shop    string        `bson:"-"`
	Buyer   string        `bson:"-"`
	BuyerId bson.ObjectId `bson:"buyer,omitempty" schema:"-"`
	ShopId  bson.ObjectId `bson:"shop,omitempty" schema:"-"`
}

type Orders struct {
	Total int       `json:"total"`
	Data  []Product `json:"data"`
}

type OrderHandler struct {
	Sess *mgo.Session
	Db   string
}

func (this *OrderHandler) Post(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	r.ParseForm()
	buyer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	order := new(Order)
	schema.NewDecoder().Decode(order, r.Form)
	if err := validator.Validate(order); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	order.Id = bson.NewObjectId()
	order.ShopId = bson.ObjectIdHex(order.Shop)
	order.Buyer = buyer
	order.State = "待付款"
	order.Date = time.Now()
	order.Price = 0
	for _, item := range order.Items {
		product := Product{}
		if err := sess.DB(this.Db).C("products").FindId(item.Id).One(&product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if product.Shop != order.ShopId {
			response.Forbidden(w, "店铺信息错误")
			return
		}
		if product.Inventory < item.Number {
			response.Forbidden(w, "库存不足")
			return
		}
		order.Price += product.Price * item.Number
	}

	if err := sess.DB(this.Db).C("orders").Insert(order); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.NoContent(w)
}

// func (this *ProductHandler) PutById(w http.ResponseWriter, r *http.Request) {
// 	sess := this.Sess.Clone()
// 	defer sess.Close()
// 	r.ParseForm()
// 	owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
// 	product := new(Product)
// 	schema.NewDecoder().Decode(product, r.Form)
// 	if err := validator.Validate(product); err != nil {
// 		response.Forbidden(w, err.Error())
// 		return
// 	}
// 	shop := Shop{}
// 	if err := sess.DB(this.Db).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	if err := sess.DB(this.Db).C("products").Update(&bson.M{"_id": bson.ObjectIdHex(r.FormValue(":id")), "shop": shop.Id}, &bson.M{"$set": product}); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	response.NoContent(w)
// }

// func (this *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
// 	sess := this.Sess.Clone()
// 	defer sess.Close()
// 	r.ParseForm()
// 	query := new(ProductQuery)
// 	schema.NewDecoder().Decode(query, r.Form)
// 	if len(query.Sort) == 0 {
// 		query.Sort = "-date"
// 	}
// 	if len(query.Word) != 0 {
// 		query.Name = map[string]string{"$regex": query.Word, "$options": "i"}
// 	}
// 	if len(query.Shop) != 0 && bson.IsObjectIdHex(query.Shop) {
// 		query.ShopId = bson.ObjectIdHex(query.Shop)
// 	}
// 	products := new(Products)
// 	if err := sess.DB(this.Db).C("products").Find(&query).Select(&bson.M{"info": 0, "props": 0, "spec": 0}).Skip(query.Skip).Limit(query.Limit).Sort("-date").All(&products.Data); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	products.Total, _ = sess.DB(this.Db).C("products").Find(&query).Count()
// 	response.JSON(w, &products)
// }

// func (this *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
// 	sess := this.Sess.Clone()
// 	defer sess.Close()
// 	id := r.FormValue(":id")
// 	if !bson.IsObjectIdHex(id) {
// 		response.Forbidden(w, "error id")
// 		return
// 	}
// 	product := Product{}
// 	if err := sess.DB(this.Db).C("products").FindId(bson.ObjectIdHex(id)).One(&product); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	response.JSON(w, &product)
// }

// func (this *ProductHandler) DelById(w http.ResponseWriter, r *http.Request) {
// 	sess := this.Sess.Clone()
// 	defer sess.Close()
// 	id, owner := r.FormValue(":id"), bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
// 	if !bson.IsObjectIdHex(id) {
// 		response.Forbidden(w, "error id")
// 		return
// 	}
// 	shop := Shop{}
// 	if err := sess.DB(this.Db).C("shops").Find(&bson.M{"owner": owner}).One(&shop); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	if err := sess.DB(this.Db).C("products").Remove(bson.M{"_id": bson.ObjectIdHex(id), "shop": shop.Id}); err != nil {
// 		response.InternalServerError(w, err.Error())
// 		return
// 	}
// 	response.NoContent(w)
// }

// func OrderRoute(m *app.App) {
// 	m.Post("/order", func(w http.ResponseWriter, r *http.Request) {
// 		buyer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
// 		data, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		var order Order
// 		if err := json.Unmarshal(data, &order); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		order.Id = bson.NewObjectId()
// 		order.Date = time.Now()
// 		order.Buyer = buyer
// 		order.State = "待付款"
// 		order.Price = 0
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		for _, item := range order.Items {
// 			product := Product{}
// 			if err := sess.DB(MongoDB).C("products").FindId(item.Id).One(&product); err != nil {
// 				response.InternalServerError(w, err.Error())
// 				return
// 			}
// 			if product.Shop != order.Shop {
// 				response.Forbidden(w, "店铺信息错误")
// 				return
// 			}
// 			if product.Inventory < item.Number {
// 				response.Forbidden(w, "库存不足")
// 				return
// 			}
// 			order.Price += product.Price * item.Number
// 		}
// 		if err := sess.DB(MongoDB).C("orders").Insert(order); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.NoContent(w)
// 	})
// }
