package store

import (
	"sync"

	"product-crud/backend/models"
)

type ProductStore struct {
	mu       sync.RWMutex
	products []models.Product
	nextID   int
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: []models.Product{
			{ID: 1, Name: "Wireless Mouse", Description: "Ergonomic wireless mouse with USB receiver", Price: 29.99, Stock: 50},
			{ID: 2, Name: "Mechanical Keyboard", Description: "RGB mechanical keyboard with blue switches", Price: 89.99, Stock: 30},
			{ID: 3, Name: "USB-C Hub", Description: "7-in-1 USB-C hub with HDMI and SD card reader", Price: 45.50, Stock: 25},
		},
		nextID: 4,
	}
}

func (s *ProductStore) GetAll() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Product, len(s.products))
	copy(result, s.products)
	return result
}

func (s *ProductStore) GetByID(id int) (models.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.products {
		if p.ID == id {
			return p, true
		}
	}
	return models.Product{}, false
}

func (s *ProductStore) Create(input models.ProductInput) models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	product := models.Product{
		ID:          s.nextID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}
	s.nextID++
	s.products = append(s.products, product)
	return product
}

func (s *ProductStore) Update(id int, input models.ProductInput) (models.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.products {
		if p.ID == id {
			s.products[i] = models.Product{
				ID:          id,
				Name:        input.Name,
				Description: input.Description,
				Price:       input.Price,
				Stock:       input.Stock,
			}
			return s.products[i], true
		}
	}
	return models.Product{}, false
}

func (s *ProductStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.products {
		if p.ID == id {
			s.products = append(s.products[:i], s.products[i+1:]...)
			return true
		}
	}
	return false
}
