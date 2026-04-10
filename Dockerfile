FROM golang:1.26-alpine3.23 as builder

WORKDIR /app
COPY . .

ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/MBI-88/*

RUN apk add --no-cache git

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build  -ldflags="-s -w" -o ./dominus cmd/api/main.go


FROM alpine:3.23.3 as deployment
WORKDIR /app

ENV REST_PORT=8000
ENV GRPC_PORT=5000
ENV MODE=prod

RUN mkdir -p /etc/dominus/certs
COPY --from=builder /app/dominus ./dominus
COPY --from=builder /app/env ./env
COPY --from=builder /app/certs /etc/dominus/certs

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown appuser:appgroup /app

RUN apk add --no-cache dos2unix && dos2unix ./env/entrypoint.sh

USER appuser
EXPOSE ${REST_PORT}
EXPOSE ${GRPC_PORT}

ENTRYPOINT ["/app/env/entrypoint.sh"]

CMD ["./dominus", "-prod=false", "-banner=true"]
