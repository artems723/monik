# Start from golang base image
FROM golang:alpine

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -ldflags "-X main.buildVersion=1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=1.0.1 " -o /build cmd/server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]