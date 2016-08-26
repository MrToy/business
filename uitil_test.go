package business

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestNewDecoder(t *testing.T) {
	s1 := &struct {
		Aa bson.ObjectId
	}{}
	hex := "57b1ae0f98fe890005000001"
	v1 := map[string][]string{"Aa": {hex}}
	decoder := NewDecoder()
	decoder.Decode(s1, v1)
	if s1.Aa != bson.ObjectIdHex(hex) {
		t.Errorf("s1.Aa: expected %v, got %v", hex, s1.Aa)
	}
}
