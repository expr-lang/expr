package playground

import "time"

type Env struct {
	Products  []Product  `expr:"products"`
	Customers []Customer `expr:"customers"`
	Discounts []Discount `expr:"discounts"`
	Orders    []Order    `expr:"orders"`
}

type Product struct {
	Name        string
	Description string
	Price       float64
	Stock       int
	AddOn       *AddOn
	Metadata    map[string]interface{}
	Tags        []string
	Rating      float64
	Reviews     []Review
}

type Feature struct {
	Id          string
	Description string
}

type Discount struct {
	Name    string
	Percent int
}

type Customer struct {
	FirstName string
	LastName  string
	Age       int
	Addresses []Address
}

type Address struct {
	Country    string
	City       string
	Street     string
	PostalCode string
}

type Order struct {
	Number    int
	Customer  Customer
	Items     []*OrderItem
	Discounts []*Discount
	CreatedAt time.Time
}

type OrderItem struct {
	Product  Product
	Quantity int
}

type Review struct {
	Product  *Product
	Customer *Customer
	Comment  string
	Rating   float64
}

type AddOn struct {
	Name  string
	Price float64
}
