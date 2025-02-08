FROM --platform=amd64 golang:latest as builder

WORKDIR /app

COPY . .
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM --platform=amd64 alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
