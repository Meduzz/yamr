FROM alpine:3.6
RUN apk update
RUN apk add postgresql ca-certificates
COPY yamr /srv/yamr/yamr
COPY static/ /srv/yamr/static/
VOLUME /srv/yamr/files
WORKDIR "/srv/yamr"
CMD ["/srv/yamr/yamr"]
EXPOSE 4040