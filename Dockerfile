# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /app

# need git for go modules
RUN apk add --no-cache git ca-certificates tzdata

# bring in module files first (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# bring in the rest of the source
COPY . .

# force excelize v2.9.0 AFTER copying repo files
RUN go get github.com/xuri/excelize/v2@v2.9.0 && go mod tidy && go mod download

# sanity check: should print Version: v2.9.0
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
