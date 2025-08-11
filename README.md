# POC LocalStack + Go (SQS)

Prova de conceito que cria uma fila SQS no LocalStack, envia uma mensagem e em seguida faz o receive + delete.

## Stack

* LocalStack (SQS) — documentação: [https://docs.localstack.cloud/](https://docs.localstack.cloud/)
* Go 1.22 — site oficial: [https://go.dev/](https://go.dev/) — instalação: [https://go.dev/doc/install](https://go.dev/doc/install)
* Aplicação Go 1.22 (rodando dentro de um container para desenvolvimento interativo via bind mount)

## Requisitos

* Docker: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
* Docker Compose: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

## Execução Rápida

Passo a passo mínimo:

```bash
# 1. Sobe os containers
docker compose up -d

# 2. Entra no container da aplicação
docker exec -it poc-localstack-app sh

# 3. Dentro do container, no diretório /app, garante dependências
go mod tidy

# 4. Executa a aplicação (envia e lê 1 mensagem)
go run .
```

Saída esperada (exemplo):

```text
2025/08/11 23:13:40 Fila pronta: demo-queue (http://poc-localstack:4566/000000000000/demo-queue)
2025/08/11 23:13:40 Mensagem enviada: Olá LocalStack! 2025-08-11T23:13:40Z
2025/08/11 23:13:40 Mensagem processada e removida: ID=0b1c334d-881a-4b5e-aca4-8de7881ad629 Body=Olá LocalStack! 2025-08-11T23:13:40Z
```

## O que a aplicação faz

A aplicação demonstra um fluxo completo de integração com SQS:

1. Configuração: Carrega variáveis de ambiente com valores padrão
2. Conexão: Estabelece conexão com LocalStack usando AWS SDK v2
3. Criação da fila: Garante que a fila `QUEUE_NAME` existe (cria se necessário)
4. Envio: Envia uma mensagem com timestamp e atributo `Origin=poc-localstack`
5. Recebimento: Faz polling da fila (máximo 1 mensagem, timeout de 2s)
6. Processamento: Exibe a mensagem recebida e a remove da fila


## Estrutura Simplificada

```text
docker-compose.yml
init-scripts/
  01-create-queue.sh        # Script executado pelo LocalStack na inicialização
app/
  main.go                   # Aplicação principal
  go.mod                    # Dependências do módulo Go
  go.sum                    # Checksums das dependências
  internal/
    env/
      env.go                # Carregamento de variáveis de ambiente
    sqs/
      config.go             # Configuração AWS/SQS
      sqs_send.go           # Envio de mensagens
      sqs_receive.go        # Recebimento e processamento de mensagens
    types/
      types.go              # Tipos compartilhados
    utils/
      helpers.go            # Funções auxiliares
  scripts/
    Dockerfile              # Dockerfile para ambiente de desenvolvimento
```

## Detalhes do LocalStack utilizados neste projeto

### Endpoints e portas

* O endpoint unificado do LocalStack é exposto na porta `4566`.
  Neste projeto, a aplicação aponta para `http://poc-localstack:4566` (hostname do serviço na rede do Docker Compose) e os comandos via CLI usam `http://localhost:4566`.
  Referência: [https://docs.localstack.cloud/aws/capabilities/networking/accessing-endpoint-url/](https://docs.localstack.cloud/aws/capabilities/networking/accessing-endpoint-url/)

### awslocal

* O `awslocal` é um wrapper sobre o `aws` já configurado para o endpoint do LocalStack e é utilizado nos comandos de exemplo e nos init-scripts.
  Repositório: [https://github.com/localstack/awscli-local](https://github.com/localstack/awscli-local)

### Init-scripts

* Os scripts montados em `/etc/localstack/init/ready.d` executam automaticamente quando o LocalStack está pronto, permitindo o provisionamento da fila via `awslocal`.
  Referência: [https://docs.localstack.cloud/aws/capabilities/config/initialization-hooks/](https://docs.localstack.cloud/aws/capabilities/config/initialization-hooks/)

## Variáveis Principais

Definidas no `docker-compose.yml`:

* `AWS_REGION` (padrão: `us-east-1`)
* `AWS_ACCESS_KEY_ID` (padrão: `test`)
* `AWS_SECRET_ACCESS_KEY` (padrão: `test`)
* `QUEUE_NAME` (padrão: `demo-queue`)
* `LOCALSTACK_ENDPOINT` (endpoint interno: `http://poc-localstack:4566`)

Os valores padrão também são aplicados no código caso as variáveis estejam ausentes.

## Dependências

* Go 1.22
* AWS SDK for Go v2

  * `github.com/aws/aws-sdk-go-v2`
  * `github.com/aws/aws-sdk-go-v2/config`
  * `github.com/aws/aws-sdk-go-v2/service/sqs`

## Comandos úteis com LocalStack

Listar todas as filas:

```bash
docker exec -it poc-localstack sh -c 'awslocal sqs list-queues'
```

Ver atributos de uma fila específica:

```bash
docker exec -it poc-localstack sh -c 'awslocal sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/demo-queue --attribute-names All'
```

Criar fila adicional manualmente:

```bash
docker exec -it poc-localstack sh -c 'awslocal sqs create-queue --queue-name outra-fila'
```

Ver mensagens na fila (sem remover):

```bash
docker exec -it poc-localstack sh -c 'awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/demo-queue'
```

## Reexecutar rapidamente

Dentro do container de app:

```bash
go run .
```

Como o código está em bind mount (`./app:/app`), alterações locais refletem imediatamente.

## Limpar ambiente

Para parar e remover todos os containers:

```bash
docker compose down
```


## Referências

* LocalStack (documentação geral): [https://docs.localstack.cloud/](https://docs.localstack.cloud/)
* Endpoints LocalStack (endpoint URL): [https://docs.localstack.cloud/aws/capabilities/networking/accessing-endpoint-url/](https://docs.localstack.cloud/aws/capabilities/networking/accessing-endpoint-url/)
* Init hooks LocalStack: [https://docs.localstack.cloud/aws/capabilities/config/initialization-hooks/](https://docs.localstack.cloud/aws/capabilities/config/initialization-hooks/)
* awscli-local (awslocal): [https://github.com/localstack/awscli-local](https://github.com/localstack/awscli-local)
* Go (site oficial): [https://go.dev/](https://go.dev/) — instalação: [https://go.dev/doc/install](https://go.dev/doc/install)
* AWS SDK for Go v2 — endpoints/exemplos SQS:
  [https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-endpoints.html](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-endpoints.html),
  [https://docs.aws.amazon.com/code-library/latest/ug/go\_2\_sqs\_code\_examples.html](https://docs.aws.amazon.com/code-library/latest/ug/go_2_sqs_code_examples.html),
  [https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs)
