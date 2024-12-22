# Build stage
FROM golang:1.23 as builder
WORKDIR /app

COPY ./go.mod /app
COPY ./go.sum /app
# RUN go mod tidy
# Copy source files
COPY ./shared /app/shared
COPY ./buyer-platform/bap-adapter-service /app/bap-adapter-service
# Download dependencies and build the binary
RUN  GOARCH=amd64 GOOS=linux go build -o app-bap-adapter-service bap-adapter-service/server.go

# Final stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/app-bap-adapter-service /app/
COPY --from=builder /app/shared/config/testdata/* /app/
RUN ls -alh /app
CMD ["./app-bap-adapter-service"]