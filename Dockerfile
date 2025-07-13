
# Build aşaması
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Çalışma aşaması
FROM alpine:3.18
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY .env .env
# ✔️ Binary dosyayı kopyala
COPY --from=builder /app/server .

# ✔️ Frontend klasörünü de kopyala
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./server"]

