package types

// Constantes de configuração padrão
const (
	DefaultRegion   = "us-east-1"
	DefaultEndpoint = "http://localhost:4566"
	DefaultQueue    = "demo-queue"
)

// EnvConfig agrega parâmetros necessários para execução.
type EnvConfig struct {
	Region    string
	Endpoint  string
	QueueName string
}

// ReceivedMessage representa o resultado de uma leitura simples.
type ReceivedMessage struct {
	ID   string
	Body string
}
