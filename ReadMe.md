
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

