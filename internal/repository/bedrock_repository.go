package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
)

// BedrockRepository handles all interactions with AWS Bedrock Agent
type BedrockRepository struct {
	client       *bedrockagentruntime.Client
	agentID      string
	agentAliasID string
}

// NewBedrockRepository creates a new BedrockRepository instance
func NewBedrockRepository(client *bedrockagentruntime.Client, agentID, agentAliasID string) *BedrockRepository {
	return &BedrockRepository{
		client:       client,
		agentID:      agentID,
		agentAliasID: agentAliasID,
	}
}

// InvokeAgent calls the Bedrock Agent and returns the streaming response
func (r *BedrockRepository) InvokeAgent(ctx context.Context, sessionID, prompt string) (*bedrockagentruntime.InvokeAgentOutput, error) {
	input := &bedrockagentruntime.InvokeAgentInput{
		AgentId:      &r.agentID,
		AgentAliasId: &r.agentAliasID,
		SessionId:    &sessionID,
		InputText:    &prompt,
		EndSession:   aws.Bool(false),
	}
	return r.client.InvokeAgent(ctx, input)
}
