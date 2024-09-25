FROM alpine:latest

RUN mkdir /app

COPY dockerServiceApp /app

CMD [ "/app/dockerServiceApp" ]