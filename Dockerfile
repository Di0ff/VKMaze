FROM golang:1.23.4 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o vkIntern

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/vkIntern .

ENTRYPOINT ["./vkIntern"]