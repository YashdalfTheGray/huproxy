FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o huproxy .

FROM scratch

WORKDIR /

COPY --from=builder /app/huproxy ./

EXPOSE 9090

CMD ["./huproxy"]
