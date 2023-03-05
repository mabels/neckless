FROM scratch

COPY ./neckless /usr/local/bin/neckless

CMD ["/usr/local/bin/neckless", "version"]
