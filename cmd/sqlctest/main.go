package main

import (
	"context"
	"math/rand"
	"sqlctest/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Db *models.Queries
}

func main() {

	context := context.Background()

	url := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

	conn, err := pgxpool.New(context, url)
	if err != nil {
		panic(err)
	}

	db := models.New(conn)
	InitDb(context, db)

	handler := Handler{Db: db}

	r := fiber.New()

	r.Get("/products", handler.AllProducts)
	r.Get("/orders", handler.AllOrders)

	r.Listen(":8080")

}

func (h *Handler) AllProducts(ctx *fiber.Ctx) error {

	products, err := h.Db.ListProducts(ctx.UserContext())

	if err != nil {
		return err
	}

	return ctx.JSON(products)

}

func (h *Handler) AllOrders(ctx *fiber.Ctx) error {

	orders, err := h.Db.ListOrders(ctx.UserContext())

	if err != nil {
		return err
	}

	return ctx.JSON(orders)

}

func InitDb(ctx context.Context, db *models.Queries) {

	count := Must(db.CountProducts(ctx))

	if count == 0 {

		c1 := Must(db.CreateCustomer(ctx, models.CreateCustomerParams{
			Name:  "John Doe",
			Email: "",
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
