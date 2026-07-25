package store

import (
	"testing"

	"product-crud/backend/models"
)

func TestNewProductStore_HasMockProducts(t *testing.T) {
	s := NewProductStore()
	products := s.GetAll()

	if len(products) != 3 {
		t.Fatalf("expected 3 mock products, got %d", len(products))
	}
	if products[0].ID != 1 || products[0].Name != "Wireless Mouse" {
		t.Fatalf("unexpected first product: %+v", products[0])
	}
}

func TestGetAll_ReturnsCopy(t *testing.T) {
	s := NewProductStore()
	products := s.GetAll()
	products[0].Name = "Modified"

	original, _ := s.GetByID(1)
	if original.Name == "Modified" {
		t.Fatal("GetAll should return a copy, not a reference to internal slice")
	}
}

func TestGetByID_Found(t *testing.T) {
	s := NewProductStore()

	product, found := s.GetByID(2)
	if !found {
		t.Fatal("expected product with id 2 to be found")
	}
	if product.Name != "Mechanical Keyboard" {
		t.Fatalf("unexpected product name: %s", product.Name)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	s := NewProductStore()

	_, found := s.GetByID(999)
	if found {
		t.Fatal("expected product with id 999 to not be found")
	}
}

func TestCreate_AutoIncrementID(t *testing.T) {
	s := NewProductStore()

	first := s.Create(models.ProductInput{
		Name:        "Laptop Stand",
		Description: "Adjustable stand",
		Price:       39.99,
		Stock:       10,
	})
	second := s.Create(models.ProductInput{
		Name:        "Webcam",
		Description: "1080p webcam",
		Price:       59.99,
		Stock:       5,
	})

	if first.ID != 4 {
		t.Fatalf("expected first new product id 4, got %d", first.ID)
	}
	if second.ID != 5 {
		t.Fatalf("expected second new product id 5, got %d", second.ID)
	}
	if first.ID == second.ID {
		t.Fatal("product IDs must be unique")
	}

	products := s.GetAll()
	if len(products) != 5 {
		t.Fatalf("expected 5 products after create, got %d", len(products))
	}
}

func TestCreate_StoresAllFields(t *testing.T) {
	s := NewProductStore()

	input := models.ProductInput{
		Name:        "Monitor",
		Description: "27 inch display",
		Price:       299.99,
		Stock:       8,
	}
	product := s.Create(input)

	if product.Name != input.Name ||
		product.Description != input.Description ||
		product.Price != input.Price ||
		product.Stock != input.Stock {
		t.Fatalf("stored product mismatch: %+v", product)
	}
}

func TestUpdate_Found(t *testing.T) {
	s := NewProductStore()

	input := models.ProductInput{
		Name:        "Wireless Mouse Pro",
		Description: "Updated mouse",
		Price:       34.99,
		Stock:       40,
	}
	updated, found := s.Update(1, input)
	if !found {
		t.Fatal("expected product id 1 to be updated")
	}
	if updated.ID != 1 {
		t.Fatalf("expected id to remain 1, got %d", updated.ID)
	}
	if updated.Name != input.Name {
		t.Fatalf("expected name %q, got %q", input.Name, updated.Name)
	}

	stored, _ := s.GetByID(1)
	if stored.Name != input.Name {
		t.Fatalf("expected persisted name %q, got %q", input.Name, stored.Name)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	s := NewProductStore()

	_, found := s.Update(999, models.ProductInput{
		Name:  "Missing",
		Price: 1,
		Stock: 1,
	})
	if found {
		t.Fatal("expected update of missing product to return false")
	}
}

func TestDelete_Found(t *testing.T) {
	s := NewProductStore()
	initialCount := len(s.GetAll())

	if !s.Delete(2) {
		t.Fatal("expected delete to succeed")
	}

	products := s.GetAll()
	if len(products) != initialCount-1 {
		t.Fatalf("expected %d products after delete, got %d", initialCount-1, len(products))
	}

	_, found := s.GetByID(2)
	if found {
		t.Fatal("deleted product should not be found")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := NewProductStore()
	initialCount := len(s.GetAll())

	if s.Delete(999) {
		t.Fatal("expected delete of missing product to return false")
	}
	if len(s.GetAll()) != initialCount {
		t.Fatalf("product count should remain %d, got %d", initialCount, len(s.GetAll()))
	}
}
