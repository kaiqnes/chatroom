FROM golang:1.20

LABEL maintainer="Caique Nunes <kaiqnes@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["go", "run", "./cmd/main.go"]
