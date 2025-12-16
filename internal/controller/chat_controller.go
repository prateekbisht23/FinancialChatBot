package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"financial-chat-bot/internal/models"
	"financial-chat-bot/internal/service"
)

// ChatController handles HTTP endpoints for chat operations
type ChatController struct {
	chatService *service.ChatService
}

// NewChatController creates a new ChatController instance
func NewChatController(chatService *service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// HandleChat handles POST /chat requests
func (c *ChatController) HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Validate HTTP method
	if r.Method != http.MethodPost {
		c.sendError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	// Get SessionID from URL Query Parameters
	sessionID := r.URL.Query().Get("sessionId")

	// Decode the request body
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(w, http.StatusBadRequest, "Invalid request format or missing prompt")
		return
	}

	// Validate required fields
	if err := c.validateRequest(sessionID, req.Prompt); err != nil {
		c.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("[Controller] Received chat request - SessionID: %s, Prompt: %s", sessionID, req.Prompt)

	// Process the chat via service layer
	response, err := c.chatService.ProcessChat(r.Context(), sessionID, req.Prompt)
	if err != nil {
		c.sendError(w, http.StatusInternalServerError, "Failed to invoke Bedrock Agent")
		return
	}

	// Handle empty response
	if response == "" {
		log.Print("[Controller] Agent returned an empty response")
		response = "Agent returned an empty response. Check agent configuration."
	}

	// Send successful response
	c.sendSuccess(w, response)
}

// validateRequest validates the incoming chat request
func (c *ChatController) validateRequest(sessionID, prompt string) error {
	if sessionID == "" {
		return &models.ValidationError{Field: "sessionId", Message: "sessionId is required in query parameters"}
	}
	if prompt == "" {
		return &models.ValidationError{Field: "prompt", Message: "prompt is required in request body"}
	}
	return nil
}

// sendError sends an error response to the client
func (c *ChatController) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ChatResponse{Error: message})
}

// sendSuccess sends a successful response to the client
func (c *ChatController) sendSuccess(w http.ResponseWriter, agentResponse string) {
	json.NewEncoder(w).Encode(models.ChatResponse{
		AgentResponse: agentResponse,
	})
}
