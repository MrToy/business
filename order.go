package business

import (
	"github.com/gohttp/response"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
	"net/http"
	"time"
)

type Order struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty" schema:"-"`
	Buyer   bson.ObjectId `json:"buyer" bson:",omitempty" schema:"-"`
	State   string        `json:"state" bson:",omitempty" schema:"-"`
	Date    time.Time     `json:"date" schema:"-"`
	Price   float64       `json:"price" schema:"-"`
	Shop    bson.ObjectId `json:"shop" bson:",omitempty"`
	Contact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Addr  string `json:"addr"`
	} `json:"contact"`
	Express struct {
		Code string `json:"code"`
		Id   string `json:"id"`
	} `json:"express"`
	Items []struct {
		Id       bson.ObjectId `json:"id"`
		Spec     string        `json:"spec"`
		Quantity int           `json:"quantity"`
		Name     string        `json:"name"`
		Price    float64       `json:"price"`
	} `json:"items"`
}

type OrderQuery struct {
	Skip   int
	Limit  int
	Sort   string
	Target string
}

type Orders struct {
	Total int     `json:"total"`
	Data  []Order `json:"data"`
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
	NewDecoder().Decode(order, r.Form)
	if err := validator.Validate(order); err != nil {
		response.Forbidden(w, err.Error())
		return
	}
	order.Id = bson.NewObjectId()
	order.Buyer = buyer
	order.State = "待付款"
	order.Date = time.Now()
	order.Price = 0
	for i, item := range order.Items {
		product := Product{}
		if err := sess.DB(this.Db).C("products").FindId(item.Id).One(&product); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if product.Shop != order.Shop {
			response.Forbidden(w, "店铺信息错误")
			return
		}
		if product.Inventory < item.Quantity {
			response.Forbidden(w, "库存不足")
			return
		}
		order.Price += product.Price * float64(item.Quantity)
		order.Items[i].Price = product.Price
		order.Items[i].Name = product.Name
	}

	if err := sess.DB(this.Db).C("orders").Insert(order); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.JSON(w, order.Id)
}

func (this *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	r.ParseForm()
	viewer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	query := new(OrderQuery)
	NewDecoder().Decode(query, r.Form)
	if len(query.Sort) == 0 {
		query.Sort = "-date"
	}
	orders := new(Orders)
	if query.Target == "buyer" {
		if err := sess.DB(this.Db).C("orders").Find(&bson.M{"buyer": viewer}).Skip(query.Skip).Limit(query.Limit).Sort("-date").All(&orders.Data); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		orders.Total, _ = sess.DB(this.Db).C("orders").Find(&bson.M{"buyer": viewer}).Count()
	} else if query.Target == "shoper" {
		shop := Shop{}
		if err := sess.DB(this.Db).C("shops").Find(bson.M{"owner": viewer}).One(&shop); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		if shop.Owner == viewer {
			if err := sess.DB(this.Db).C("orders").Find(&bson.M{"shop": shop.Id}).Skip(query.Skip).Limit(query.Limit).Sort("-date").All(&orders.Data); err != nil {
				response.InternalServerError(w, err.Error())
				return
			}
			orders.Total, _ = sess.DB(this.Db).C("orders").Find(&bson.M{"shop": shop.Id}).Count()
		}
	}
	response.JSON(w, &orders)
}

func (this *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := r.FormValue(":id")
	if !bson.IsObjectIdHex(id) {
		response.Forbidden(w, "error id")
		return
	}
	order := Order{}
	if err := sess.DB(this.Db).C("orders").FindId(bson.ObjectIdHex(id)).One(&order); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.JSON(w, &order)
}

func (this *OrderHandler) DeliverById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := bson.ObjectIdHex(r.FormValue(":id"))
	viewer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	order := Order{}
	if err := sess.DB(this.Db).C("orders").FindId(id).One(&order); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	shop := Shop{}
	if err := sess.DB(this.Db).C("shops").Find(bson.M{"owner": viewer}).One(&shop); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if shop.Id != order.Shop {
		response.Forbidden(w, "没有足够的权限")
		return
	}
	if err := sess.DB(this.Db).C("orders").UpdateId(id, &bson.M{"$set": &bson.M{"express": &bson.M{"id": r.FormValue("id"), "code": r.FormValue("code")}}}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if order.State == "待发货" {
		if err := sess.DB(this.Db).C("orders").UpdateId(id, &bson.M{"$set": &bson.M{"state": "已发货"}}); err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
	}
}

func (this *OrderHandler) CancleById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := bson.ObjectIdHex(r.FormValue(":id"))
	viewer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	if err := sess.DB(this.Db).C("orders").Update(&bson.M{"state": "待付款", "_id": id, "buyer": viewer}, &bson.M{"$set": &bson.M{"state": "已取消"}}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
}

func (this *OrderHandler) ConfirmById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := bson.ObjectIdHex(r.FormValue(":id"))
	viewer := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
	if err := sess.DB(this.Db).C("orders").Update(&bson.M{"state": "已发货", "_id": id, "buyer": viewer}, &bson.M{"$set": &bson.M{"state": "交易完成"}}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
}
