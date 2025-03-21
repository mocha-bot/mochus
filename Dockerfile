FROM golang:1.22.12-alpine3.21 as builder

RUN apk update && apk upgrade && \
  apk --no-cache --update add git make

WORKDIR /app

COPY . .

RUN go mod tidy && \
  go mod download && \
  go build -v -o engine && \
  chmod +x engine

## Distribution
FROM alpine:latest

# Install dependencies
RUN apk update && apk upgrade && \
  apk --no-cache --update add ca-certificates tzdata && \
  mkdir mochus

# Install Doppler CLI
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler

WORKDIR /mochus

COPY --from=builder /app/engine /mochus
