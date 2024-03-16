# Use a multi-stage build for the Go application
FROM golang:1.19.0

WORKDIR /usr/src/app

# Copy only the go.mod and go.sum files to optimize caching
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /usr/src/app/main .

# Set the binary as executable
RUN chmod 755 /usr/src/app/main

# Expose the application port
EXPOSE 3000
EXPOSE 5432

# Command to run the application
CMD ["/usr/src/app/main"]
