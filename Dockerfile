FROM golang:1.19.2-bullseye as builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 go build -o nory ./cmd/server

ARG DATABASE_CERT_URL
RUN curl --create-dirs -o /root.crt -O ${DATABASE_CERT_URL}

FROM gcr.io/distroless/static-debian11 as production

COPY --from=builder /app/nory /
COPY --from=builder /root.crt /.postgresql/root.crt

CMD ["/nory"]
