# Projeto Cotacao USD/BRL (Go + SQLite)

Este projeto expõe um endpoint HTTP que consulta a cotacao USD/BRL em uma API externa, grava o resultado em SQLite e retorna apenas o valor do `bid`. Um cliente simples faz chamadas ao servidor para demonstrar timeouts.

## Arquitetura (resumo)

- `cmd/server`: servidor HTTP.
- `cmd/client`: cliente HTTP de exemplo.
- `pkg/types`: tipos compartilhados (DTOs).
- `app.db`: banco SQLite criado automaticamente.

## Endpoints

- `GET /cotacao`

Resposta:

```json
{ "bid": "5.1234" }
```

## Como rodar

Requer Go 1.25.6.

Build:

```bash
make build
```

Rodar servidor:

```bash
make run-server
```

Rodar cliente:

```bash
make run-client
```

Rodar servidor + cliente (cliente dispara e o servidor encerra depois):

```bash
make run
```

## Configuracoes importantes

No servidor (`cmd/server/main.go`):

- `serverAddr`: `:8080`
- `priceAPIURL`: `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- `callTimeout`: 200ms (timeout da chamada externa)
- `dbTimeout`: 10ms (timeout para inserir no SQLite)

No cliente (`cmd/client/main.go`):

- `apiAddress`: `http://localhost:8080/cotacao`
- `callTimeout`: 300ms (chamadas que devem funcionar)
- `callTimeoutToFail`: 2ms (chamadas que devem falhar)
- `apiCalls`: 5 (numero de chamadas em cada lote)

## Banco de dados

O servidor cria a tabela `prices` ao iniciar e limpa a tabela a cada start. O arquivo do banco e `app.db` na raiz do projeto.

## Observaçoes

- A resposta do endpoint retorna apenas o campo `bid`.
- Erros de chamada externa ou escrita no banco retornam `500`.
