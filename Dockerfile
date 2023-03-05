FROM scratch

COPY ./neckless /bin/neckless

ENTRYPOINT ["neckless"]
