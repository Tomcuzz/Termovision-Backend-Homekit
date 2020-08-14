# Start from the latest golang base image
FROM golang:latest

# Set the mantainer
MAINTAINER "Thomas Cousin"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/tomcuzz/Termovision-Backend-Homekit/src

ENV HK_PIN="00102003"
ENV HK_SERIAL="027TC-000001"

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src .

# Inatall dependancies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ./main -hkpin ${HK_PIN} -hkserial ${HK_SERIAL}
