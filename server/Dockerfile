# syntax=docker/dockerfile:1

FROM golang:1.24-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /server

EXPOSE 5199

CMD ["/server"]
