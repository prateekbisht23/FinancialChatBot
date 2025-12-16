package service

import (
	"context"
	"log"

	"financial-chat-bot/internal/repository"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

// ChatService handles the business logic for chat operations
type ChatService struct {
	bedrockRepo *repository.BedrockRepository
}

// NewChatService creates a new ChatService instance
func NewChatService(bedrockRepo *repository.BedrockRepository) *ChatService {
	return &ChatService{
		bedrockRepo: bedrockRepo,
	}
}

// ProcessChat processes a chat request and returns the agent's response
func (s *ChatService) ProcessChat(ctx context.Context, sessionID, prompt string) (string, error) {
	log.Printf("[Service] Processing chat for SessionID: %s", sessionID)

	// Invoke the Bedrock Agent via repository
	resp, err := s.bedrockRepo.InvokeAgent(ctx, sessionID, prompt)
	if err != nil {
		log.Printf("[Service] Bedrock Agent Invoke Error: %v", err)
		return "", err
	}

	// Extract response from stream
	response := s.extractResponseFromStream(resp)
	return response, nil
}

// extractResponseFromStream reads all chunks from the event stream
func (s *ChatService) extractResponseFromStream(resp *bedrockagentruntime.InvokeAgentOutput) string {
	var finalResponse string

	stream := resp.GetStream()
	if stream == nil {
		return ""
	}
	defer stream.Close()

	// Iterate over the events channel
	for event := range stream.Events() {
		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			if len(e.Value.Bytes) > 0 {
				finalResponse += string(e.Value.Bytes)
			}
		case *types.ResponseStreamMemberTrace:
			// Log trace events for debugging (optional)
			// log.Printf("[Service] Trace event received")
			continue
		default:
			continue
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		log.Printf("[Service] Stream error: %v", err)
	}

	return finalResponse
}
