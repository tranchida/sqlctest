package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sqlctest/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Db *models.Queries
}

func main() {

	context := context.Background()

	var url string
	if url = os.Getenv("POSTGRESQL_URL"); url == "" {
		url = "postgres://gouser:password@localhost:5432/mydb?sslmode=disable"
	}

	conn, err := pgxpool.New(context, url)
	if err != nil {
		panic(err)
	}

	db := models.New(conn)

	gin.SetMode(gin.ReleaseMode)
	log.Printf("GIN server starting version : %s\n", gin.Version)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handler := Handler{Db: db}

	r.GET("/customers", handler.AllCustomers)
	r.GET("/products", handler.AllProducts)
	r.GET("/orders", handler.AllOrders)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server started on port 8080")
	log.Fatal(r.Run(":8080"))

}

func (h *Handler) AllCustomers(c *gin.Context) {

	customers, err := h.Db.ListCustomers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)

}

func (h *Handler) AllProducts(c *gin.Context) {

	products, err := h.Db.ListProducts(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)

}

func (h *Handler) AllOrders(c *gin.Context) {

	orders, err := h.Db.ListOrders(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}
