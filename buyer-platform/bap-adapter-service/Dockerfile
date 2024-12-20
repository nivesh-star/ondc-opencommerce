# Build stage
FROM golang:1.23 as builder
WORKDIR /app

# Copy source files
COPY ./bap-adapter-service /app/bap-adapter-service
COPY ./shared /app/shared

# Set the working directory for package `a`
WORKDIR /app/bap-adapter-service

# Download dependencies and build the binary
RUN go mod tidy && go build -o bap-adapter-service .

# Final stage
FROM alpine:3.17
WORKDIR /root/
COPY --from=builder /app/bap-adapter-service/bap-adapter-service .
CMD ["./bap-adapter-service"]
