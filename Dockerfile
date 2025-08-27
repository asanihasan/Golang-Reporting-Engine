# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /app

# need git for go modules; certs & tzdata are nice to have
RUN apk add --no-cache git ca-certificates tzdata

# bring in module files first (best cache)
COPY go.mod go.sum ./

# if your repo's go.mod still pins excelize 2.8.x, force 2.10.0 here
# (delete these two lines if you've already updated go.mod in GitHub)
RUN go get github.com/xuri/excelize/v2@v2.10.0 && go mod tidy

# fetch deps
RUN go mod download

# bring in source
COPY . .

# sanity check the resolved excelize version
RUN go list -m -json github.com/xuri/excelize/v2 | grep Version

# build statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# ---------- runtime stage ----------
FROM alpine:3.20
WORKDIR /app

# folders your app writes to
RUN mkdir -p /app/source /app/result \
  && addgroup -S app && adduser -S app -G app \
  && chown -R app:app /app
USER app

COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 6969
ENV GIN_MODE=release
CMD ["server"]
