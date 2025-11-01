# Movie API

Uma API RESTful para gerenciar uma biblioteca de filmes com Go. A aplicação fornece endpoints para criar, ler e deletar registros de filmes com health check para deploys containerizados.
## Tabela de Conteúdo
- [Features](#features)
- [Pré-requisitos](#pré-requisitos)
- [Build](#build)
- [Desenvolvimento](#desenvolvimento)
- [API](#api)
- [Estrutura](#estrutura)
- [Testes](#testes)
- [Requisitos](#requisitos)
- [Configuração](#configuração)
## Features
- RESTful API com respostas JSON
- Health check para orquestração de containeres
- Containerizado com Docker
- Makefile para implementação com um clique (one-click deployment) e tarefas de desenvolvimento
- Comunicação interna utilizando Protocol Buffers gRPC
- Testes mockados e e2e
- Documentação via Swagger
- Versionamento de rotas `/v1` visando evolução da API sem perda de compatibilidade
## Pré-requisitos
### Para rodar a aplicação
- **Docker** ou **Podman** (containerização)
- **Make** (para automação de building)
### Para desenvolvimento
- **Go** 1.25+ (Go development)
- **Protocol Buffer Compiler** (suporte para protobuf)
- **Make** (para automação de build e tarefas de desenvolvimento)

## Build
### Usando Docker/Podman
1. **Clone o repositório**
   ```bash
   git clone https://github.com/Kahlis/teste-tecnico-sipub-tech.git
   cd teste-tecnico-sipub-tech
   ```

2. **Build do container**
   ```bash
   # With Docker (or Podman via alias)
   make build
   ```
Esse comando cria uma cópia do .env.example com um ambiente já configurado de produção para fins de testes

3. **Testar a aplicação funcionando**
   ```bash
   curl http://localhost:8080/v1/health
   ```

A API estará disponível via bridge em `http://localhost:8080`
### Lista de comandos Make
```bash
# Gera o .env e sobe o docker compose
make build

# Gera as definições Protobuf para movies e para apigateway
make protoc

# Sobe o docker compose
make up

# Desliga o docker compose
make down

# Roda a bateria de testes mockados
make mock

# Roda a bateria de testes em produção dentro do container (deve ser rodado)
make e2e

# Apaga os arquivos de armazenamento do container mongodb
make clean

# Instala, verifica e linka todas as dependências do Go
make deps
```
## Desenvolvimento
### Configurando o ambiente de desenvolvimento
1. **Instala, verifica e linka as dependências Go**
   ```bash
   make deps
   ```
2. **Gera os Protocol Buffers**
   ```bash
   make proto
   ```
3. **Roda testes mockados para verificar a integridade**
   ```bash
   make mock
   ```
## API
Para acessar a documentação da API, acesse [http://localhost:8080/v1/swagger/index.html](http://localhost:8080/swagger/index.html)
### Formato de resposta
Todas as respostas seguem o seguinte formato:
```json
{
  "status": "success|error",
  "data": {}
}
```
### Endpoints
#### Health Check
```bash
# Check service health
curl http://localhost:8080/v1/health

# Response
{"status":"success","data":null}
```
#### List Movies
```bash
# Lista todos os filmes com paginação
curl http://localhost:8080/v1/movies

# Resposta
{
  "status": "success",
  "data": {
    "movies": [
      {"id": "1", "title": "Inception", "year": "2010"},
      ...
    ],
    "more": false,
    "page": 1,
    "total": 2,
    "results": 2
  }
}
```
#### Get Movie by ID
```bash
# Retorna um filme pelo Id
curl http://localhost:8080/v1/movies/1

# Resposta
{
  "status": "success",
  "data": {
    "id": "1",
    "title": "Inception",
    "year": "2010"
  }
}
```
#### Create Movie
```bash
# Registra um novo filme
curl -X POST http://localhost:8080/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Interstellar", "year": "2014"}'

# Resposta
{
  "status": "success",
  "data": {
    "id": "3",
    "title": "Interstellar",
    "year": "2014"
  }
}
```
#### Delete Movie
```bash
# Exlui um filme pelo Id
curl -X DELETE http://localhost:8080/v1/movies/1

# Resposta: 204 No Content
```
## Estrutura
### apigateway
```
apigateway/
├── go.mod                 # Módulo Go
├── Dockerfile             # Definição do container docker
├── cmd/
│   └── api/
│       └── main.go        # Entry point da aplicação API
├── core/
│   ├── app/               # Configuração do app via dependency injection
│   ├── config/            # Configuração do ambiente .env
│   ├── routes/            # Registrador de rotas
│   ├── domain/            # Camada que armazena as regras de negócio
│   ├── usecases/          # Camada de implementação do domínio
│   ├── handlers/          # Camada de handlers HTTP
│   ├── proto/             # Arquivos Protobuf auto gerados
│   └── util/              # Armazenamento de utilidades e erros
├── infra/clients/         # Definição dos clientes e seus contratos
├── pkg/logger/            # Configuração de um logger centralizado
└── routes/                # API route definitions
```
### movies
```
movies/
├── go.mod                 # Módulo Go
├── Dockerfile             # Definição do container docker
├── cmd/
│   └── server/
│       └── main.go        # Entry point da aplicação Movies
├── core/
│   ├── app/               # Configuração do app via dependency injection
│   ├── config/            # Configuração do ambiente .env
│   ├── integration/       # Camada de integrações
│   │   └──  test/         # Armazenamento de testes
│   │        └──  e2e/     # Testes e2e em ambiente simulado de produção
│   │        └──  mock/    # Testes mockados
│   ├── domain/            # Camada que armazena as regras de negócio
│   ├── usecases/          # Camada de implementação do domínio
│   ├── handlers/          # Camada de handlers HTTP
│   ├── repository/        # Camada de definição da interface com o db
│   ├── proto/             # Arquivos Protobuf auto gerados
│   └── util/              # Armazenamento de utilidades e erros
├── infra/persistence/     # Definição dos clientes e seus contratos
├── pkg/logger/            # Configuração de um logger centralizado
├── seed/                  # Armazenamento dos dados iniciais do banco
└── routes/                # API route definitions
```
### root
```
teste-tecnico-sipub-tech/
├── .gitignore             # .gitignore padrão Go
├── docker-compose.yaml    # Compose para orquestração do container
├── Makefile               # Comandos de desenvolvimento e build
├── REQUIREMENTS.md        # Levantamento de requisitos inicial
├── README.md              # Informações sobre o projeto
├── Makefile               # Comandos de desenvolvimento e build
├── apigateway/
│   └── .../               # Arquivos do microsserviço apigateway
├── movies/
│   └── .../               # Arquivos do microsserviço movies
└── proto/                 # Arquivos de definição .proto
```
## Testes
O projeto possui dois tipos principais de testes: **E2E (end-to-end)** e **mockados**, cada um com seu propósito.

### Testes E2E
```bash
make e2e
````

Os testes **E2E** são executados dentro do container **Movies** e verificam a integração completa do **API Gateway**.
> Execute esses testes na primeira vez que inicializar o repositório.
> Caso já tenha containers antigos, rode:
>
> ```bash
> make clean
> make build
> make e2e
> ```
>
> Isso garante um ambiente limpo e uma execução correta dos testes de ponta a ponta.


### Testes Mockados
```bash
make mock
```

Os testes **mockados** rodam diretamente no código, sem a necessidade de usar o `docker compose`.
Eles utilizam um repositório simulado para validar a lógica interna e a estrutura da aplicação, garantindo o funcionamento adequado das camadas de negócio e suas integrações.
## Requisitos
### Dependências Go
Gerenciado via `go.mod` - `make deps` verifica e baixa todas as dependências para cada microsserviço.
### Dependências do Container
- Base image: `golang:1.25.3-alpine`
- Port: `8080` exposed
- Health check: `GET /v1/health`
## Configuração
A aplicação pode ser configurada utilizando os valores no arquivo `.env`:

```bash
ENV=
API_PORT=
MOVIES_PORT=
MONGO_CONTAINER_NAME=
MONGO_DB=
MONGO_DB_USER=
MONGO_DB_PASSWORD=
MONGO_DB_COLLECTION=
MONGO_DB_URI=
```
