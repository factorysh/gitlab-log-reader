FROM bearstech/python-dev:3.9 as dev

RUN mkdir -p /opt/logger \
        && python3 -m venv /opt/logger/venv \
        && /opt/logger/venv/bin/pip install -U pip wheel \
        && /opt/logger/venv/bin/pip install flask

COPY fake_log.py /opt/logger/

FROM bearstech/python:3.9

COPY --from=dev /opt/logger /opt/logger
ARG uid=500

USER ${uid}
CMD [ "/opt/logger/venv/bin/python", "/opt/logger/fake_log.py" ]
