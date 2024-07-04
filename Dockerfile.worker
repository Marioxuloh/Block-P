# Use the official Golang image as a base
FROM golang:1.22.1

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application to the container
COPY . .

# Build the Go app
RUN go build -o app

# Expose the port on which the Go app will run
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
