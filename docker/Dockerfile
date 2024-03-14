FROM golang:1.22

# Set destination for COPY
WORKDIR /code

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

ADD . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /guildmaster ./cmd

EXPOSE 8080

# Run
CMD ["/guildmaster"]