package models

type Data struct {
	Order    Orders   `json:"order"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Items  `json:"items"`
}
