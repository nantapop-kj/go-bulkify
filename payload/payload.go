package payload

import "fmt"

func Wrap[T any](fn func(index int) (T, string)) func(int) (any, string) {
	return func(index int) (any, string) {
		return fn(index)
	}
}

type Payload struct {
	ProductName  string `json:"product_name"`
	ProductCode  string `json:"product_code"`
	DisplayOrder int    `json:"display_order"`
}

func BuildPayload(index int) (any, string) {
	name := fmt.Sprintf("Product_%05d", index)
	return Payload{
		ProductName:  name,
		ProductCode:  fmt.Sprintf("SKU-%05d", index),
		DisplayOrder: index,
	}, name
}
