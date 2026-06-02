FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/floatip ./cmd/floatip/

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D floatip
COPY --from=builder /app/floatip /usr/local/bin/floatip
USER floatip
ENTRYPOINT ["floatip"]
