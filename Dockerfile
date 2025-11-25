FROM golang:1.25 as builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

RUN go mod download

RUN go build -o /opt/go/transactions -v cmd/main.go
RUN chmod -R +x /opt/go/transactions

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder --chown=nonroot:nonroot /opt/go /opt/go
COPY --from=builder --chown=nonroot:nonroot /app/migrations migrations

USER nonroot
ENTRYPOINT ["/opt/go/transactions"]
