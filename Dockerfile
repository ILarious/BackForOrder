FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/back-for-order ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/back-for-order-worker ./cmd/worker

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/back-for-order /app/back-for-order
COPY --from=builder /bin/back-for-order-worker /app/back-for-order-worker

EXPOSE 8080

CMD ["/app/back-for-order"]
