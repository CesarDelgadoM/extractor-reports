package server

import (
	"github.com/CesarDelgadoM/extractor-reports/config"
	"github.com/CesarDelgadoM/extractor-reports/internal/extractors/branch"
	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/internal/requests"
	"github.com/CesarDelgadoM/extractor-reports/internal/workerpool"
	"github.com/CesarDelgadoM/extractor-reports/pkg/database"
	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Run() {
	// Logger
	zap.InitLogger(s.config)

	// Database
	mongodb := database.ConnectMongoDB(s.config.Mongo)
	defer mongodb.Disconnect()

	// Stream
	rabbitmq := stream.ConnectRabbitMQ(s.config.RabbitMQ)
	defer rabbitmq.Close()

	// Workerpool
	workerpool := workerpool.NewWorkerPool(s.config.Worker)

	// Store
	store := requests.NewStoreRequests()

	// Producer
	opts := &producer.ProducerOpts{
		ExchangeName: "reports",
		ExchangeType: "direct",
		ContentType:  "application/json",
	}
	producer := producer.NewProducer(opts, rabbitmq)
	defer producer.Close()

	// App
	app := fiber.New()

	// Repositorys
	branchRepository := branch.NewBranchRepository(mongodb)

	// Extractors
	branchExtractor := branch.NewBranchExtractor(store, producer, branchRepository)

	// Services
	branchService := branch.NewBranchService(workerpool, store, branchExtractor)

	// Handlers
	branch.NewBranchHandler(app, branchService)

	// Launch
	app.Listen(s.config.Server.Port)
}
