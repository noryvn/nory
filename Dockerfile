FROM golang:1.19.2-bullseye as builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 go build -o nory ./cmd/server

FROM debian:bullseye as production

COPY --from=builder /app/nory /

RUN curl --create-dirs -o $HOME/.postgresql/root.crt -O ${DATABASE_CERT_URL}

CMD ["/nory"]
