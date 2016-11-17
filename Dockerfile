FROM alpine:3.4
RUN apk update
RUN apk add postgresql
COPY yamr /opt/yamr/
VOLUME /opt/yamr/files
ENTRYPOINT ["/opt/yamr/yamr"]
EXPOSE 4040