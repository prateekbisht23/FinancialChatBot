package router

import (
	"encoding/json"
	"net/http"

	"financial-chat-bot/internal/controller"
)

// Router handles all route registrations
type Router struct {
	chatController *controller.ChatController
}

// NewRouter creates a new Router instance
func NewRouter(chatController *controller.ChatController) *Router {
	return &Router{
		chatController: chatController,
	}
}

// RegisterRoutes registers all HTTP routes
func (router *Router) RegisterRoutes() {
	http.HandleFunc("/chat", router.chatController.HandleChat)
	http.HandleFunc("/health", router.healthCheck)
}

// healthCheck is a simple health check endpoint
func (router *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "FinancialChatBot",
	})
}
