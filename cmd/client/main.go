package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventory "Study/gen"
)

// GatewayServer - REST шлюз для доступа к gRPC сервису инвентаря
type GatewayServer struct {
	inventoryClient inventory.InventoryServiceClient
}

// NewGatewayServer - создает новый экземпляр REST шлюза
func NewGatewayServer() *GatewayServer {
	return &GatewayServer{}
}

// ConnectToGRPC - подключается к gRPC сервису
func (g *GatewayServer) ConnectToGRPC(address string) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	g.inventoryClient = inventory.NewInventoryServiceClient(conn)
	return nil
}

// StockResponse - структура ответа с количеством товара
type StockResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

// GetStockHandler - обработчик HTTP запроса для получения количества товара
func (g *GatewayServer) GetStockHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем product_id из параметра маршрута
	vars := mux.Vars(r)
	productID := vars["id"]
	if productID == "" {
		http.Error(w, "Необходим параметр id в маршруте", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	response, err := g.inventoryClient.GetStock(ctx, &inventory.StockRequest{
		ProductId: productID,
	})
	if err != nil {
		slog.Error("Ошибка при получении данных от gRPC сервиса", "error", err, "product_id", productID)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StockResponse{
		ProductID: productID,
		Quantity:  response.GetQuantity(),
	})

	slog.Info("Получено количество товара через шлюз", "product_id", productID, "quantity", response.GetQuantity())
}

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Создание REST шлюза
	gateway := NewGatewayServer()

	// Подключение к gRPC сервису
	// В Docker используем имя сервиса, локально - localhost
	grpcAddress := "inventory-server:50051" // Имя сервиса в docker-compose
	// grpcAddress := "localhost:50051" // Для локального запуска
	if err := gateway.ConnectToGRPC(grpcAddress); err != nil {
		slog.Error("Ошибка подключения к gRPC сервису", "error", err, "address", grpcAddress)
		os.Exit(1)
	}

	// Настройка HTTP маршрутов с использованием gorilla/mux
	router := mux.NewRouter()
	router.HandleFunc("/stock/{id}", gateway.GetStockHandler).Methods("GET")

	// Создание HTTP сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	slog.Info("REST шлюз запущен", "address", server.Addr)

	// Graceful shutdown - правильная остановка сервера
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("Остановка REST шлюза...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Ошибка при остановке шлюза", "error", err)
		} else {
			slog.Info("REST шлюз остановлен")
		}
	}()

	// Запуск сервера
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Ошибка REST шлюза", "error", err)
		os.Exit(1)
	}
}
