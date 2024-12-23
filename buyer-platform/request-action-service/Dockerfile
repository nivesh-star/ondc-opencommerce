# Build stage
FROM golang:1.23 as builder
WORKDIR /app

COPY ./go.mod /app
COPY ./go.sum /app
# RUN go mod tidy
# Copy source files
COPY ./shared /app/shared
COPY ./buyer-platform/request-action-service /app/request-action-service
# Download dependencies and build the binary
ARG TARGETOS
ARG TARGETARCH
RUN  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o app-request-action-service request-action-service/server.go

# Final stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/app-request-action-service /app/
COPY --from=builder /app/shared/config/testdata/* /app/
RUN ls -alh /app
CMD ["./app-request-action-service"]