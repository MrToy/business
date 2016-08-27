package business

import (
	//"fmt"
	//"github.com/gohttp/app"
	//"github.com/gohttp/response"
	//"labix.org/v2/mgo"
	"gopkg.in/mgo.v2/bson"
	//"net/http"
	//"strconv"
	//"time"
)

type Contact struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Owner bson.ObjectId `json:"owner"`
	Name  string        `json:"name"`
	Phone string        `json:"phone"`
	Addr  string        `json:"addr"`
}

// func ContactRoute(m *app.App) {
// 	m.Post("/contacts", func(w http.ResponseWriter, r *http.Request) {
// 		owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
// 		name, phone, addr := r.FormValue("name"), r.FormValue("phone"), r.FormValue("addr")
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		if err := sess.DB(MongoDB).C("contacts").Insert(&Contact{Id: bson.NewObjectId(), Owner: owner, Name: name, Phone: phone, Addr: addr}); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.NoContent(w)
// 	})
// 	m.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
// 		owner := bson.ObjectIdHex(r.Header.Get("X-Auth-Id"))
// 		contacts := []Contact{}
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		if err := sess.DB(MongoDB).C("contacts").Find(&bson.M{"owner": owner}).All(&contacts); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.JSON(w, contacts)
// 	})
// 	m.Get("/contacts/:id", func(w http.ResponseWriter, r *http.Request) {
// 		id := r.FormValue(":id")
// 		if !bson.IsObjectIdHex(id) {
// 			response.Forbidden(w, "error id")
// 			return
// 		}
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		res := Contact{}
// 		if err := sess.DB(MongoDB).C("contacts").FindId(bson.ObjectIdHex(id)).One(&res); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.JSON(w, res)
// 	})
// 	m.Del("/contacts/:id", func(w http.ResponseWriter, r *http.Request) {
// 		id, owner := r.FormValue(":id"), r.Header.Get("X-Auth-Id")
// 		if !bson.IsObjectIdHex(id) {
// 			response.Forbidden(w, "error id")
// 			return
// 		}
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		if err := sess.DB(MongoDB).C("contacts").Remove(bson.M{"_id": bson.ObjectIdHex(id), "owner": bson.ObjectIdHex(owner)}); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.NoContent(w)
// 	})
// 	m.Put("/contacts/:id", func(w http.ResponseWriter, r *http.Request) {
// 		id, owner := r.FormValue(":id"), r.Header.Get("X-Auth-Account")
// 		if !bson.IsObjectIdHex(id) {
// 			response.Forbidden(w, "error id")
// 			return
// 		}
// 		sess := Sess.Clone()
// 		defer sess.Close()
// 		data := bson.M{}
// 		if name := r.FormValue("name"); len(name) != 0 {
// 			data["name"] = name
// 		}
// 		if phone := r.FormValue("phone"); len(phone) != 0 {
// 			data["phone"] = phone
// 		}
// 		if addr := r.FormValue("addr"); len(addr) != 0 {
// 			data["addr"] = addr
// 		}
// 		if err := sess.DB(MongoDB).C("products").Update(bson.M{"_id": bson.ObjectIdHex(id), "owner": owner}, bson.M{"$set": data}); err != nil {
// 			response.InternalServerError(w, err.Error())
// 			return
// 		}
// 		response.NoContent(w)
// 	})
// }
