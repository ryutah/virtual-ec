package model

type ProductID int

type Product struct {
	id    ProductID
	name  string
	price int
}

func NewProduct(id ProductID, name string, price int) *Product {
	return &Product{
		id:    id,
		name:  name,
		price: price,
	}
}

func (p *Product) ID() ProductID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() int {
	return p.price
}
