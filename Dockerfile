ARG ARCH

##
# BUILD CONTAINER
##

FROM alpine:3.12 as builder

RUN \
apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM ${ARCH}/busybox:1.32-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY approuvez /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/approuvez"]
CMD [""]
