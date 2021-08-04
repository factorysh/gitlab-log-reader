FROM bearstech/debian:bullseye

ARG uid=500
COPY bin/log-reader-ip /usr/local/bin/

USER ${uid}
EXPOSE 8000
EXPOSE 8042
CMD [ "/usr/local/bin/log-reader-ip" ]
