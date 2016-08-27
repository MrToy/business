package business

import (
	"github.com/gohttp/response"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type PayHandler struct {
	Sess *mgo.Session
	Db   string
}
type NotifyQuery struct {
	Out_trade_no bson.ObjectId `schema:"out_trade_no"`
	Trade_status string        `schema:"trade_status"`
}

func (this *PayHandler) RedirectById(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	id := r.FormValue(":id")
	order := Order{}
	if err := sess.DB(this.Db).C("orders").FindId(bson.ObjectIdHex(id)).One(&order); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	if order.State != "待付款" {
		response.Forbidden(w, "订单状态错误")
		return
	}
	urlStr, _ := GetPayUrl(order.Id.Hex(), order.Items[0].Name, order.Price, r.FormValue("notify"))
	http.Redirect(w, r, urlStr, http.StatusFound)
}

func (this *PayHandler) Notify(w http.ResponseWriter, r *http.Request) {
	sess := this.Sess.Clone()
	defer sess.Close()
	query := new(NotifyQuery)
	r.ParseForm()
	NewDecoder().Decode(query, r.Form)
	if !VerifyPaySign(r.Form) {
		log.Printf("订单 %s 验证信息错误", query.Out_trade_no.Hex())
		return
	}
	if query.Trade_status != "TRADE_FINISHED" && query.Trade_status != "TRADE_SUCCESS" {
		log.Printf("订单 %s 状态错误：%s", query.Out_trade_no.Hex(), query.Trade_status)
		return
	}
	if err := sess.DB(this.Db).C("orders").UpdateId(query.Out_trade_no, &bson.M{"$set": &bson.M{"state": "待发货"}}); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	log.Printf("订单 %s 支付成功", query.Out_trade_no.Hex())
	response.OK(w, "success")
}
