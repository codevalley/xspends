# Start from the official Go image which contains all the necessary build tools and libraries
FROM golang:1.21 AS build

# Set the current working directory inside the container
WORKDIR /app

ENV DB_DSN="root:@tcp(tidb-cluster-tidb.tidb-cluster.svc.cluster.local:4000)/xspends"

# Copy the local package files to the container's workspace
COPY . .

# Build the application
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o xpends .

# Use a smaller image to run the application
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary file from the build container to the production container
COPY --from=build /app/xpends .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./xpends"]
