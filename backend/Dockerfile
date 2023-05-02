FROM golang:1.20-alpine as builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/main.go

FROM alpine

# Install ca-certificates for SSL/TLS support
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Expose the server port
EXPOSE 8080

# Set the command to run the program with the EXTERNAL_RNG environment variable as an argument
CMD ["./main", "--log-type", "json", "--addr", ":8080", "--rng", "$EXTERNAL_RNG"]
