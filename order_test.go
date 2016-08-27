package business

import (
	"testing"
)

func TestOrderSchema(t *testing.T) {
	order := new(Order)
	v1 := map[string][]string{
		"items.0.name":     {"衣服"},
		"items.0.price":    {"100"},
		"items.0.quantity": {"1"},
		"items.0.spec":     {"颜色：黄,尺寸：XL"},
		"items.0.id":       {"57ba92a419a18b0005e58435"},
		"shop":             {"57b1ae0f98fe890005000001"},
	}
	NewDecoder().Decode(order, v1)
	t.Log(order)
}

func TestGetPayUrl(t *testing.T) {
	str, _ := GetPayUrl("123", "test", 1.22, "http://")
	t.Log(str)
}
