# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# bring in module files first for cache
COPY go.mod go.sum ./
RUN go mod download

# now copy the rest of the source
COPY . .

# FORCE excelize v2.9.0 *after* repo files are in place
RUN go get github.com/xuri/excelize/v2@v2.9.0 && go mod tidy && go mod download

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
