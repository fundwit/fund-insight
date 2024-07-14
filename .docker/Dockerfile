FROM busybox:1.36.1-glibc

COPY fund-insight /
COPY buildInfo.json /
ENV GIN_MODE release
ENTRYPOINT ["/fund-insight"]