# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# 1) copy module files first to leverage cache
COPY go.mod go.sum ./

# 2) force the exact excelize version + tidy
#    (this updates go.mod/go.sum *inside the image*)
RUN go get github.com/xuri/excelize/v2@v2.9.0 && go mod tidy
RUN go mod download

# 3) now copy the rest of the source and build
COPY . .
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
