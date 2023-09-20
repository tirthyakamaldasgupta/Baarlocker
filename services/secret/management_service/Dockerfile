FROM golang:latest

RUN apt-get update && apt-get install -y uuid-runtime

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
