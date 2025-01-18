package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sqlctest/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	Db *models.Queries
}

func main() {

	context := context.Background()
	url := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

	conn, err := pgx.Connect(context, url)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context)

	Db := models.New(conn)
	InitDb(context, Db)

	handler := Handler{Db: Db}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/products", handler.AllProducts)

	r.Run()

}

func (h *Handler) AllProducts(ctx *gin.Context) {

	orders, err := h.Db.ListOrders(ctx)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})

}

func InitDb(ctx context.Context, Db *models.Queries) {

	count, err := Db.CountProducts(ctx)
	if err != nil {
		panic(err)
	}

	if count == 0 {


		c1 := Must(Db.CreateCustomer(ctx, models.CreateCustomerParams{
			Name:  "John Doe",
			Email: "",
		}))

		p1 := Must(Db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample1",
			Price: 30,
			Stock: int32(rand.Intn(100)),
		}))

		p2 := Must(Db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample2",
			Price: 20,
			Stock: int32(rand.Intn(100)),
		}))

		Must(Db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample3",
			Price: 30,
			Stock: int32(rand.Intn(100)),
		}))

		Must(Db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample4",
			Price: 40,
			Stock: int32(rand.Intn(100)),
		}))
		
		Must(Db.CreateProduct(ctx, models.CreateProductParams{
			Code:  "sample5",
			Price: 50,
			Stock: int32(rand.Intn(100)),
		}))


		Must(Db.CreateOrder(ctx, models.CreateOrderParams{
			CustomerID: c1.ID,
			ProductID:  p1.ID,
			Quantity:   1,
		}))

		Must(Db.CreateOrder(ctx, models.CreateOrderParams{
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
