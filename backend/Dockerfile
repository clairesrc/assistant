FROM --platform=amd64 golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./src ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM --platform=amd64 alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
