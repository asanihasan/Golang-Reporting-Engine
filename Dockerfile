# ---------- build stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app

# install certs (https for go mod) + tzdata (nice to have)
RUN apk add --no-cache ca-certificates tzdata

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# ---------- runtime stage ----------
FROM alpine:3.20
WORKDIR /app

# create the folders your code expects (source/ and result/)
RUN mkdir -p /app/source /app/result \
    && addgroup -S app && adduser -S app -G app \
    && chown -R app:app /app
USER app

COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 6969
ENV GIN_MODE=release
CMD ["server"]
