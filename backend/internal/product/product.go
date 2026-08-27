package product

import "octacart/internal/shared"

type Product struct {
	id    int
	name  string
	sku   string
	price shared.Money
	attrs map[string]string //update this according to spec i provided.
}

// use getters with entities like following style. read the spec before you touch any go file.

func (p *Product) Id() int             { return p.id }
func (p *Product) Name() string        { return p.name }
func (p *Product) Sku() string         { return p.sku }
func (p *Product) Price() shared.Money { return p.price }
