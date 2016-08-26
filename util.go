package business

import (
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

func NewDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(bson.NewObjectId(), func(s string) reflect.Value {
		return reflect.ValueOf(bson.ObjectIdHex(s))
	})
	return decoder
}
