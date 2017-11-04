FROM alpine:3.4
RUN apk update
RUN apk add postgresql
COPY yamr /srv/yamr/yamr
COPY static/ /srv/yamr/static/
VOLUME /srv/yamr/files
VOLUME /srv/yamr/gce
WORKDIR "/srv/yamr"
CMD ["/srv/yamr/yamr"]
EXPOSE 4040