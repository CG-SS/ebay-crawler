package model

type ItemModel struct {
	Id         string `json:"-"`
	Title      string `json:"title"`
	Condition  string `json:"condition"`
	Price      string `json:"price"`
	ProductUrl string `json:"product_url"`
}
