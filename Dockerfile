FROM golang:latest


WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY .env .env

RUN go build -o /app/main ./cmd/otp

EXPOSE 8090
CMD ["/app/main"]