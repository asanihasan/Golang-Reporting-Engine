# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# bring in module files first (best cache)
COPY go.mod go.sum ./
RUN go mod download

# bring in the rest of the source
COPY . .

# sanity check which excelize got resolved (should show v2.10.0)
RUN go list -m -json github.com/xuri/excelize/v2 | grep Version

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# ---------- runtime stage ----------
FROM alpine:3.20
WORKDIR /app
RUN mkdir -p /app/source /app/result \
    && addgroup -S app && adduser -S app -G app \
    && chown -R app:app /app
USER app

COPY --from=builder /app/server /usr/local/bin/server
EXPOSE 6969
ENV GIN_MODE=release
CMD ["server"]
