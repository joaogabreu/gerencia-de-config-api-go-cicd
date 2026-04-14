# API em Go

## Sobre o projeto
API de CRUD de usuários para prática de Docker, Compose e CI/CD.

## O que foi implementado
- Persistência em PostgreSQL (com criação automática da tabela `users`)
- `Dockerfile` multi-stage para build e execução da API
- `docker-compose.yml` com API + PostgreSQL + PGAdmin
- Testes unitários e de integração
- Pipeline CI/CD em GitHub Actions cobrindo CI, container e CD

## Pré-requisitos
- Go >= 1.23
- Docker e Docker Compose

## Rodando local sem Docker
1. Entrar na pasta da API:
   - `cd "APIs para exercícios/api-go"`
2. Baixar dependências:
   - `go mod tidy`
3. Subir PostgreSQL local ou usar variáveis de ambiente para uma instância existente
4. Executar aplicação:
   - `go run main.go`

Variáveis suportadas:
- `DATABASE_URL`
- `DB_HOST` (padrão: `localhost`)
- `DB_PORT` (padrão: `5432`)
- `DB_NAME` (padrão: `api_go`)
- `DB_USER` (padrão: `postgres`)
- `DB_PASSWORD` (padrão: `postgres`)
- `DB_SSLMODE` (padrão: `disable`)

## Atividade 1: Docker
Build da imagem:
- `docker build -t api-go .`

Run do container:
- `docker run --rm -p 3000:3000 --name api-go api-go`

## Atividade 2: Docker Compose
Subir API + PostgreSQL + PGAdmin:
- `docker compose up -d --build`

Serviços:
- API: `http://localhost:3010`
- PGAdmin: `http://localhost:8081`
  - Email: `admin@admin.com`
  - Senha: `admin`

## Desafio: CI/CD
Workflow criado em:
- `.github/workflows/api-go-cicd.yml`

Etapas implementadas:
- CI
  - Build
  - Testes unitários
  - Testes de integração (com PostgreSQL em service container)
  - Lint com `golangci-lint`
  - SonarQube
  - SAST com Semgrep
- Container
  - Docker Lint com Hadolint
  - Build da imagem
  - Scan de vulnerabilidade com Trivy
  - Push da imagem no Docker Hub
- CD
  - Deploy de homolog no Render
  - DAST com OWASP ZAP
  - Aprovação manual via environment `production-approval`
  - Deploy de produção no Render

## Secrets necessários no GitHub
- `DOCKERHUB_USER`
- `DOCKERHUB_PWD`
- `SONAR_TOKEN`
- `SONAR_HOST_URL`
- `SONAR_PROJECT_KEY`
- `SONAR_ORGANIZATION`
- `RENDER_DEPLOY_HOOK_HOMOLOG`
- `HOMOLOG_APP_URL`
- `RENDER_DEPLOY_HOOK_PROD`

## Observação importante
Para a aprovação manual funcionar, configure o ambiente `production-approval` no GitHub com reviewers obrigatórios.