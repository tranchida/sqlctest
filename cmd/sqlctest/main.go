package main

import (
	"context"
	"log"
	"math/rand"
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
	InitDb(context, db)

	gin.SetMode(gin.ReleaseMode)
	log.Printf("GIN server starting version : %s\n", gin.Version)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handler := Handler{Db: db}

	r.GET("/customers", handler.AllCustomers)
	r.GET("/products", handler.AllProducts)
	r.GET("/orders", handler.AllOrders)

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

func InitDb(ctx context.Context, db *models.Queries) {

	count := Must(db.CountProducts(ctx))

	if count == 0 {

		c1 := Must(db.CreateCustomer(ctx, models.CreateCustomerParams{
			Name:  "John Doe",
			Email: "jdoe@fake.net",
		}))

		p1 := Must(db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample1",
			Price: 30,
			Stock: int32(rand.Intn(100)),
		}))

		p2 := Must(db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample2",
			Price: 20,
			Stock: int32(rand.Intn(100)),
		}))

		Must(db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample3",
			Price: 30,
			Stock: int32(rand.Intn(100)),
		}))

		Must(db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample4",
			Price: 40,
			Stock: int32(rand.Intn(100)),
		}))

		Must(db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample5",
			Price: 50,
			Stock: int32(rand.Intn(100)),
		}))

		Must(db.CreateOrder(ctx, models.CreateOrderParams{
			CustomerID: c1.ID,
			ProductID:  p1.ID,
			Quantity:   1,
		}))

		Must(db.CreateOrder(ctx, models.CreateOrderParams{
			CustomerID: c1.ID,
			ProductID:  p2.ID,
			Quantity:   2,
		}))

	}

}

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
