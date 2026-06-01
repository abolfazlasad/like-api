ARG GO_VERSION
ARG ALPINE_VERSION

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
COPY docs/ ./docs/
COPY internal/ ./internal/
RUN go build -o server ./main.go

FROM alpine:${ALPINE_VERSION}
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
