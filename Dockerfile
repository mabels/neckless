FROM ubuntu:latest

COPY ./neckless-linux /usr/local/bin/neckless

CMD ["/usr/local/bin/neckless", "version"]