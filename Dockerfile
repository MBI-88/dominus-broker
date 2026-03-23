FROM golang:1.26.1-alpine3.23.3 as builder

WORKDIR /app
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build cmd/api/main.go -ldflags="-s -w" -o ./dominus


FROM alpine:3.23.3 as deployment
WORKDIR /app

ENV REST_PORT=8000
ENV GRPC_PORT=5000

COPY --from=builder /app/dominus ./dominus
RUN mkdir -p /var/dominus/logs
RUN mkdir -p /etc/dominus/certs
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown appuser:appgroup /app

USER appuser
EXPOSE ${REST_PORT}
EXPOSE ${GRPC_PORT}

ENTRYPOINT [ "./env/entrypoint.sh" ]
CMD ["sh","-c", "./dominus -prod=true -banner=true"]
