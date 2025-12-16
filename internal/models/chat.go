package models

// ChatRequest represents the incoming chat request from the client
type ChatRequest struct {
	Prompt string `json:"prompt"`
}

// ChatResponse contains the agent's response
type ChatResponse struct {
	AgentResponse string `json:"agentResponse"`
	Error         string `json:"error,omitempty"`
}
