FROM alpine:latest
WORKDIR /app
COPY commit-stream .
ENTRYPOINT [ "/app/commit-stream" ]