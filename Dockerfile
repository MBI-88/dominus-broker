FROM golang:1.24.4-alpine3.22 as builder

WORKDIR /app
COPY . .
RUN go mod tidy && \ 
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./dominus


FROM alpine:3.20 as deployment 
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

CMD ["sh","-c", "./dominus -prod=true -banner=false"]