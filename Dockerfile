FROM golang:1.23.2-alpine

# Creating work directory   
WORKDIR /app

# Copying all files from the current directory to the container
COPY . .

# Downloading dependencies
RUN go mod tidy
RUN go build -o loan-engine ./cmd/main.go

# Command to run the executable
CMD ["./loan-engine"]