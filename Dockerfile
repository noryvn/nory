FROM golang:1.19.2-bullseye as builder

WORKDIR /app

COPY go.mod go.sum /app
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 go build -o nory ./cmd/server

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/nory /

CMD ["/nory"]
