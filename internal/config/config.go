package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AgentID      string
	AgentAliasID string
	ServerAddr   string
}

func LoadConfig() *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	agentID := os.Getenv("BEDROCK_AGENT_ID")
	agentAliasID := os.Getenv("BEDROCK_AGENT_ALIAS_ID")
	if agentID == "" || agentAliasID == "" {
		log.Fatal("BEDROCK_AGENT_ID or BEDROCK_AGENT_ALIAS_ID not set")
	}

	return &Config{
		AgentID:      agentID,
		AgentAliasID: agentAliasID,
		ServerAddr:   "192.168.1.61:8080", // Using the specific IP as requested in previous code
	}
}
