
# Test Bank Ina

This repository is for testing purpose


## Deployment

To deploy this repo we can refer to make file.

## Run Locally

Clone the project

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd my-project
```

Run docker compose for database

```bash
  docker compose up -d
```

Run tracer and metric server if you want

```bash
  cd tracing && docker compose up -d
```

Run main file

```bash
  go run main.go
```

## Swagger Open API

Install Swag/Swaggo

```bash
  go install github.com/swaggo/swag/cmd/swag@latest
```

Update Swag docs

```bash
  swag init
```

Run the application and open swagger UI

```bash
  http://localhost:8080/swagger/index.html
```