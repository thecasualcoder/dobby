FROM dobby

COPY ./bootstrap /usr/local/bin
ENTRYPOINT ["/usr/local/bin/bootstrap"]
