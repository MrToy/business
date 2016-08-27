package business

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func NewDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(bson.NewObjectId(), func(s string) reflect.Value {
		return reflect.ValueOf(bson.ObjectIdHex(s))
	})
	return decoder
}

func GetPayUrl(id, name string, price float64, notify string) (string, error) {
	keyData, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return "", errors.New("public key error")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	query := url.Values{
		"partner":        {"2088221651662605"},
		"service":        {"create_direct_pay_by_user"},
		"_input_charset": {"utf-8"},
		"timestamp":      {time.Now().Format("2006-01-02 15:04:05")},
		"version":        {"1.0"},
		"out_trade_no":   {id},
		"total_fee":      {strconv.FormatFloat(price, 'f', 2, 32)},
		"subject":        {name},
		"payment_type":   {"1"},
		"seller_email":   {"1203111636@qq.com"},
		"notify_url":     {notify},
	}
	str, _ := url.QueryUnescape(query.Encode())
	hashed := sha1.Sum([]byte(str))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}
	crypted := base64.StdEncoding.EncodeToString(signature)
	query["sign"] = []string{crypted}
	query["sign_type"] = []string{"RSA"}
	queryurl, _ := url.Parse("https://mapi.alipay.com/gateway.do")
	queryurl.RawQuery = query.Encode()
	return queryurl.String(), nil
}

func VerifyPaySign(query url.Values) bool {
	keyData, _ := ioutil.ReadFile("public.pem")
	block, _ := pem.Decode(keyData)
	pub, _ := x509.ParsePKIXPublicKey(block.Bytes)
	rsaPub, _ := pub.(*rsa.PublicKey)
	sign, sign_type := query.Get("sign"), query.Get("sign_type")
	if sign_type != "RSA" {
		return false
	}
	delete(query, "sign")
	delete(query, "sign_type")
	str, _ := url.QueryUnescape(query.Encode())
	hashed := sha1.Sum([]byte(str))
	data, _ := base64.StdEncoding.DecodeString(sign)
	if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, hashed[:], data); err != nil {
		return false
	}
	return true
}
