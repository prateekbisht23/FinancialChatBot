package main

import (
	"context"
	"log"
	"net/http"

	"financial-chat-bot/internal/config"
	"financial-chat-bot/internal/controller"
	"financial-chat-bot/internal/repository"
	"financial-chat-bot/internal/router"
	"financial-chat-bot/internal/service"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
)

// App holds all application dependencies
type App struct {
	bedrockRepo    *repository.BedrockRepository
	chatService    *service.ChatService
	chatController *controller.ChatController
	router         *router.Router
}

// NewApp initializes the application with all dependencies
func NewApp(cfg *config.Config) *App {
	// Initialize AWS client
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	bedrockClient := bedrockagentruntime.NewFromConfig(awsCfg)

	// Initialize layers (Dependency Injection)
	bedrockRepo := repository.NewBedrockRepository(bedrockClient, cfg.AgentID, cfg.AgentAliasID)
	chatService := service.NewChatService(bedrockRepo)
	chatController := controller.NewChatController(chatService)
	router := router.NewRouter(chatController)

	return &App{
		bedrockRepo:    bedrockRepo,
		chatService:    chatService,
		chatController: chatController,
		router:         router,
	}
}

// Run starts the HTTP server
func (a *App) Run(addr string) error {
	a.router.RegisterRoutes()
	log.Printf("Server starting on %s...", addr)
	log.Printf("Access from other devices: http://%s/chat", addr)
	return http.ListenAndServe(addr, nil)
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	app := NewApp(cfg)

	if err := app.Run(cfg.ServerAddr); err != nil {
		log.Fatal(err)
	}
}
