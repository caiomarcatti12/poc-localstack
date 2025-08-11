package env

import (
	"os"

	"github.com/caiomarcatti12/poc-localstack/internal/types"
)

// LoadEnvConfig lê variáveis de ambiente e aplica defaults.
func LoadEnvConfig() types.EnvConfig {
	ensureDefaultEnv("AWS_REGION", types.DefaultRegion)
	ensureDefaultEnv("AWS_ACCESS_KEY_ID", "test")
	ensureDefaultEnv("AWS_SECRET_ACCESS_KEY", "test")

	endpoint := firstNonEmpty(os.Getenv("LOCALSTACK_ENDPOINT"), types.DefaultEndpoint)
	queue := firstNonEmpty(os.Getenv("QUEUE_NAME"), types.DefaultQueue)

	return types.EnvConfig{
		Region:    os.Getenv("AWS_REGION"),
		Endpoint:  endpoint,
		QueueName: queue,
	}
}

// ensureDefaultEnv define valor padrão se variável estiver vazia.
func ensureDefaultEnv(key, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}

// firstNonEmpty retorna a primeira string não vazia.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
