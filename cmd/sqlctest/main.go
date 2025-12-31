package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sqlctest/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	Db *models.Queries
}

func main() {

	context := context.Background()

	var url string
	if url = os.Getenv("POSTGRESQL_URL"); url == "" {
		url = "postgres://user:password@localhost:5432/sqlctest?sslmode=disable"
	}

	conn, err := pgxpool.New(context, url)
	if err != nil {
		panic(err)
	}

	db := models.New(conn)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("URI: %s, status: %d", v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	handler := Handler{Db: db}

	e.GET("/customers", handler.AllCustomers)
	e.GET("/products", handler.AllProducts)
	e.GET("/orders", handler.AllOrders)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/seed", handler.Seed)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(e.Start(":8080"))

}

func (h *Handler) AllCustomers(c echo.Context) error {

	customers, err := h.Db.ListCustomers(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, customers)

}

func (h *Handler) AllProducts(c echo.Context) error {

	products, err := h.Db.ListProducts(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)

}

func (h *Handler) AllOrders(c echo.Context) error {

	orders, err := h.Db.ListOrders(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)

}

func (h *Handler) Seed(c echo.Context) error {

	customer, err := h.Db.CreateCustomer(context.Background(), models.CreateCustomerParams{Name: "John", Email: "john@fake.com"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, customer)

}
