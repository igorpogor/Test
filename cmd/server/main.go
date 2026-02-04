package main

import (
	inventory "Study/gen"
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// InventoryServer - сервер инвентаря с встроенной базой данных
type InventoryServer struct {
	inventory.UnimplementedInventoryServiceServer
	mu      sync.RWMutex     // Мьютекс для защиты от race condition
	stockDB map[string]int32 // Встроенная база данных (product_id -> количество)
}

// NewInventoryServer - создает новый экземпляр сервера инвентаря
func NewInventoryServer() *InventoryServer {
	return &InventoryServer{
		stockDB: make(map[string]int32),
	}
}

// GetStock - получает количество товара по его ID
func (s *InventoryServer) GetStock(ctx context.Context, req *inventory.StockRequest) (*inventory.StockResponse, error) {
	s.mu.RLock()         // Блокировка для чтения
	defer s.mu.RUnlock() // Разблокировка при выходе

	quantity, exists := s.stockDB[req.GetProductId()]
	if !exists {
		slog.Warn("Товар не найден", "product_id", req.GetProductId())
		return &inventory.StockResponse{Quantity: 0}, nil
	}

	slog.Info("Получено количество товара", "product_id", req.GetProductId(), "quantity", quantity)
	return &inventory.StockResponse{Quantity: quantity}, nil
}

// AddStock - добавляет товар в базу данных
func (s *InventoryServer) AddStock(productID string, quantity int32) {
	s.mu.Lock()         // Блокировка для записи
	defer s.mu.Unlock() // Разблокировка при выходе
	s.stockDB[productID] = quantity
	slog.Info("Товар добавлен в базу", "product_id", productID, "quantity", quantity)
}

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Создание сервера
	server := NewInventoryServer()

	// Добавление начальных данных в базу
	server.AddStock("product1", 100)
	server.AddStock("product2", 50)
	server.AddStock("product3", 75)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	inventory.RegisterInventoryServiceServer(grpcServer, server)

	// Включение рефлексии для отладки
	reflection.Register(grpcServer)

	// Запуск прослушивания порта
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Ошибка при запуске сервера", "error", err)
		os.Exit(1)
	}

	slog.Info("gRPC сервер запущен", "address", listener.Addr().String())

	// Graceful shutdown - правильная остановка сервера
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("Остановка сервера...")

		grpcServer.GracefulStop()
		slog.Info("Сервер остановлен")
	}()

	// Запуск сервера
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Ошибка сервера", "error", err)
		os.Exit(1)
	}
}
