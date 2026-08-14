FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/blog-server ./cmd/server

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl tzdata \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/blog-server ./blog-server
COPY --from=builder /src/internal/static ./internal/static

RUN mkdir -p /app/data/uploads

EXPOSE 8888

CMD ["./blog-server"]
