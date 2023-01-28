# Start from golang base image
FROM golang:alpine

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /build cmd/server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]