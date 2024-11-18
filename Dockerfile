# Use the official Golang image as the base image
FROM golang:latest

# Install necessary tools
RUN apt-get update && \
    apt-get install -y btrfs-progs

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .
# Ensure all files in /app are owned by root
RUN chown -R root:root /app
# Start Generation Here
RUN go mod tidy
# End Generation Here

# Run the tests
#CMD ["go", "test", "-v", "cmd/create_test.go"]
CMD ["go", "test", "./..."]
