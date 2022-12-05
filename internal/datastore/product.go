package datastore

type (
	Product struct {
		ID   uint64
		Name string
	}

	ProductFilter struct {
		IDs   []uint64
		Names []string
	}
)

type ProductRepo interface {
	CRUD[Product, ProductFilter]
}
