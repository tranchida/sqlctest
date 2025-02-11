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
		url = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	}

	conn, err := pgxpool.New(context, url)
	if err != nil {
		panic(err)
	}

	err = initDB(conn)
	if err != nil {
		log.Println("table already exist")
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

func initDB(d *pgxpool.Pool) error {

	// Lire et exécuter les fichiers SQL

	sqlBytes, err := os.ReadFile("build/sqlc/schema.sql")
	if err != nil {
		return err
	}

	_, err = d.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		return err
	}

	log.Printf("Fichier SQL exécuté avec succès")
	return nil

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
