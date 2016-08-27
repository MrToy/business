package business

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
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

func TestVerifyPaySign(t *testing.T) {
	str := "buyer_email=1659808224%40qq.com&buyer_id=2088412225884631&discount=0.00&gmt_create=2016-08-27+16%3A09%3A44&gmt_payment=2016-08-27+16%3A09%3A47&is_total_fee_adjust=N&notify_id=6087283ab8ffba9e32d535377737306kv2&notify_time=2016-08-27+16%3A13%3A36&notify_type=trade_status_sync&out_trade_no=57c14abe28313e0005e51f34&payment_type=1&price=0.01&quantity=1&seller_email=1203111636%40qq.com&seller_id=2088221651662605&sign=BI0V4YIYbbdLNwj1gd5HRZiI6Wc%2B5JgM8GH7c3Ur%2BQIbbUu7I5P6%2F6zT5jPMbT0hoGit58%2BoEBNK%2FqjxPIno0Xt%2BUPueLHTVy1XEESfrl1BEjF4g15nkp7%2FJx3xE3A4cCnex4gJRUdsAs7WfuL1nwQ3imljeaXoOx%2F33zPkzNzM%3D&sign_type=RSA&subject=test&total_fee=0.01&trade_no=2016082721001004630265360024&trade_status=TRADE_SUCCESS&use_coupon=N"
	query, _ := url.ParseQuery(str)
	if !VerifyPaySign(query) {
		t.Error("can't verify")
	}
}
