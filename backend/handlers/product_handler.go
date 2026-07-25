package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"product-crud/backend/models"
	"product-crud/backend/store"
)

type ProductHandler struct {
	store *store.ProductStore
}

func NewProductHandler(s *store.ProductStore) *ProductHandler {
	return &ProductHandler{store: s}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products := h.store.GetAll()
	return c.JSON(products)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid product ID"})
	}

	product, found := h.store.GetByID(id)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Product not found"})
	}

	return c.JSON(product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var input models.ProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid JSON request body"})
	}

	input.Name = strings.TrimSpace(input.Name)

	if errMsg := validateProductInput(input); errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: errMsg})
	}

	product := h.store.Create(input)
	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid product ID"})
	}

	var input models.ProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid JSON request body"})
	}

	input.Name = strings.TrimSpace(input.Name)

	if errMsg := validateProductInput(input); errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: errMsg})
	}

	product, found := h.store.Update(id, input)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Product not found"})
	}

	return c.JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid product ID"})
	}

	if !h.store.Delete(id) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Product not found"})
	}

	return c.JSON(models.MessageResponse{Message: "Product deleted successfully"})
}

func parseID(idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid product id")
	}
	return id, nil
}

func validateProductInput(input models.ProductInput) string {
	if strings.TrimSpace(input.Name) == "" {
		return "Name is required and cannot be empty"
	}
	if input.Price < 0 {
		return "Price must be greater than or equal to 0"
	}
	if input.Stock < 0 {
		return "Stock must be greater than or equal to 0"
	}
	return ""
}
