# Multi Stage Docker build to reduce size
# First stage
FROM golang:1.22.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Second stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/certs ./certs
COPY --from=builder /app/.env  .env
EXPOSE 3000
CMD ["./main"]
