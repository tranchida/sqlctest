package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sqlctest/internal/models"

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

	http.HandleFunc("GET /products", handler.AllProducts)
	http.HandleFunc("GET /orders", handler.AllOrders)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func (h *Handler) AllProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.Db.ListProducts(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Error writing response: %v", err)
	}

}

func (h *Handler) AllOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := h.Db.ListOrders(context.Background())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Error writing response: %v", err)
	}

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
