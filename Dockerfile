FROM alpine:latest
WORKDIR /app
COPY commit-stream .
CMD [ ./commit-stream ]