FROM golang:1.20.3-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
CMD ["go", "run", "proxy/main.go", "2000"]
