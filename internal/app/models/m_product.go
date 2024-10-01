package models

type ProductInfo struct {
	ID       int
	Serial   string `title:"Serial"`
	GsCode   string
	Model    string
	Category string
}
