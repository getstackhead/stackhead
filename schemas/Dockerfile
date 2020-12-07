FROM alpine:3.12

WORKDIR /var/www
VOLUME /var/www
COPY ./*.schema.json ./html/
COPY ./examples html/examples
COPY ./project-definition html/project-definition

CMD tail -f /dev/null
