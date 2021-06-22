FROM bearstech/debian:bullseye

ARG uid=500
COPY bin/log-reader-ip /usr/local/bin/

USER ${uid}
CMD [ "/usr/local/bin/log-reader-ip" ]
