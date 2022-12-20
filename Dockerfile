FROM golang:1.18 As builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o image-store cmd/image-store/main.go

FROM alpine:latest As production
WORKDIR /app
COPY --from=builder /app/image-store .

COPY .env .
COPY internal/db/migration ./db/migration

# executable
CMD [ "/app/image-store" ]