package models

import "time"

type Order struct {
	Id          int       `json:"id"`
	ProductId   int       `json:"product_id"`
	UserId      int       `json:"user_id"`
	Quantity    int       `json:"quantity"`
	DateCreated time.Time `json:"date_created"`
}

type RequestOrder struct {
	ProductId int `json:"product_id"`
	UserId    int `json:"user_id"`
	Quantity  int `json:"quantity"`
}

type DetailOrder struct {
	Name        string `json:"name"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
