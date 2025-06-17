
# Build the application
go build -o golang_cron .

# Create a Dockerfile
FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o golang_cron

EXPOSE 8080

CMD ["./golang_cron"]