FROM golang:1.27.2-alpine AS builder

WORKDIR /build

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o skill-badges ./cmd/skill-badges

FROM alpine:3.24

WORKDIR /app

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

COPY --chown=appuser:appgroup --from=builder /build/skill-badges .
COPY --chown=appuser:appgroup --from=builder /build/assets ./assets

EXPOSE 8080

ENTRYPOINT [ "./skill-badges" ]
