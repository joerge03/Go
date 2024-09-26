FROM alpine:latest

RUN mkdir /app

COPY mailServiceApp /app
COPY template /template

CMD [ "/app/mailServiceApp" ]