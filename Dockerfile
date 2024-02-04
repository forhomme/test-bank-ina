FROM golang:1.21-alpine as builder
LABEL stage=builder

RUN apk add --no-cache make git build-base
ENV CGO_ENABLED=0

COPY .  /usr/app/go
WORKDIR /usr/app/go
RUN go env
RUN go build -ldflags="-s -w" -o test-bank-ina

FROM ubuntu:22.04
WORKDIR /app

COPY --from=builder /usr/app/go/test-bank-ina /app/
COPY --from=builder /usr/app/go/config/config.yaml /app/config
COPY --from=builder /usr/app/go/migration/ /app/migration
COPY --from=builder /usr/app/go/docs/ /app/docs

RUN apt-get update && apt-get upgrade -y && apt-get install tzdata ca-certificates -y
ENV TZ="Asia/Singapore"
EXPOSE 8081

CMD ["./test-bank-ina"]
