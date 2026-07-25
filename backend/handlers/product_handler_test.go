package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"product-crud/backend/models"
	"product-crud/backend/store"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	handler := NewProductHandler(store.NewProductStore())

	products := app.Group("/api/products")
	products.Get("/", handler.GetAll)
	products.Get("/:id", handler.GetByID)
	products.Post("/", handler.Create)
	products.Put("/:id", handler.Update)
	products.Delete("/:id", handler.Delete)

	return app
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body io.Reader) (*http.Response, []byte) {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return resp, respBody
}

func TestGetAll_ReturnsProducts(t *testing.T) {
	app := setupTestApp()

	resp, body := doRequest(t, app, fiber.MethodGet, "/api/products/", nil)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var products []models.Product
	if err := json.Unmarshal(body, &products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(products) < 3 {
		t.Fatalf("expected at least 3 products, got %d", len(products))
	}
}

func TestGetByID_Found(t *testing.T) {
	app := setupTestApp()

	resp, body := doRequest(t, app, fiber.MethodGet, "/api/products/1", nil)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if product.ID != 1 || product.Name != "Wireless Mouse" {
		t.Fatalf("unexpected product: %+v", product)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	app := setupTestApp()

	resp, body := doRequest(t, app, fiber.MethodGet, "/api/products/999", nil)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}

	var errResp models.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error != "Product not found" {
		t.Fatalf("unexpected error message: %q", errResp.Error)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	app := setupTestApp()

	resp, body := doRequest(t, app, fiber.MethodGet, "/api/products/abc", nil)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	var errResp models.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error != "Invalid product ID" {
		t.Fatalf("unexpected error message: %q", errResp.Error)
	}
}

func TestCreate_Success(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Laptop Stand","description":"Adjustable stand","price":39.99,"stock":15}`
	resp, body := doRequest(t, app, fiber.MethodPost, "/api/products/", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", resp.StatusCode, string(body))
	}

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if product.ID == 0 || product.Name != "Laptop Stand" {
		t.Fatalf("unexpected created product: %+v", product)
	}
}

func TestCreate_EmptyName(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"   ","description":"No name","price":10,"stock":1}`
	resp, body := doRequest(t, app, fiber.MethodPost, "/api/products/", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	var errResp models.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error != "Name is required and cannot be empty" {
		t.Fatalf("unexpected error message: %q", errResp.Error)
	}
}

func TestCreate_NegativePrice(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Invalid Price","description":"Bad price","price":-1,"stock":1}`
	resp, _ := doRequest(t, app, fiber.MethodPost, "/api/products/", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreate_NegativeStock(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Invalid Stock","description":"Bad stock","price":10,"stock":-1}`
	resp, _ := doRequest(t, app, fiber.MethodPost, "/api/products/", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	app := setupTestApp()

	resp, _ := doRequest(t, app, fiber.MethodPost, "/api/products/", bytes.NewBufferString(`{invalid`))
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdate_Success(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Wireless Mouse Pro","description":"Updated","price":34.99,"stock":45}`
	resp, body := doRequest(t, app, fiber.MethodPut, "/api/products/1", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.StatusCode, string(body))
	}

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if product.ID != 1 || product.Name != "Wireless Mouse Pro" {
		t.Fatalf("unexpected updated product: %+v", product)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Missing","description":"Not found","price":10,"stock":1}`
	resp, _ := doRequest(t, app, fiber.MethodPut, "/api/products/999", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	app := setupTestApp()

	payload := `{"name":"Bad ID","description":"Invalid","price":10,"stock":1}`
	resp, _ := doRequest(t, app, fiber.MethodPut, "/api/products/0", bytes.NewBufferString(payload))
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestDelete_Success(t *testing.T) {
	app := setupTestApp()

	resp, body := doRequest(t, app, fiber.MethodDelete, "/api/products/1", nil)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var msg models.MessageResponse
	if err := json.Unmarshal(body, &msg); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if msg.Message != "Product deleted successfully" {
		t.Fatalf("unexpected message: %q", msg.Message)
	}

	getResp, _ := doRequest(t, app, fiber.MethodGet, "/api/products/1", nil)
	if getResp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected deleted product to return 404, got %d", getResp.StatusCode)
	}
}

func TestDelete_NotFound(t *testing.T) {
	app := setupTestApp()

	resp, _ := doRequest(t, app, fiber.MethodDelete, "/api/products/999", nil)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		input   string
		wantID  int
		wantErr bool
	}{
		{"1", 1, false},
		{"42", 42, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		id, err := parseID(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("parseID(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("parseID(%q) unexpected error: %v", tt.input, err)
		}
		if id != tt.wantID {
			t.Fatalf("parseID(%q) = %d, want %d", tt.input, id, tt.wantID)
		}
	}
}

func TestValidateProductInput(t *testing.T) {
	tests := []struct {
		name    string
		input   models.ProductInput
		wantErr string
	}{
		{
			name: "valid input",
			input: models.ProductInput{
				Name: "Valid", Price: 0, Stock: 0,
			},
			wantErr: "",
		},
		{
			name:    "empty name",
			input:   models.ProductInput{Name: "   ", Price: 10, Stock: 1},
			wantErr: "Name is required and cannot be empty",
		},
		{
			name:    "negative price",
			input:   models.ProductInput{Name: "Item", Price: -0.01, Stock: 1},
			wantErr: "Price must be greater than or equal to 0",
		},
		{
			name:    "negative stock",
			input:   models.ProductInput{Name: "Item", Price: 10, Stock: -1},
			wantErr: "Stock must be greater than or equal to 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateProductInput(tt.input)
			if got != tt.wantErr {
				t.Fatalf("validateProductInput() = %q, want %q", got, tt.wantErr)
			}
		})
	}
}
