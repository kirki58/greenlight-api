FROM postgres:16-alpine

COPY ./initdb/ /docker-entrypoint-initdb.d/

RUN chmod 500 /docker-entrypoint-initdb.d && \
    find /docker-entrypoint-initdb.d/ -type f -name "*.sh" -exec chmod 500 {} +

RUN chown postgres:postgres /docker-entrypoint-initdb.d && \
    find /docker-entrypoint-initdb.d/ -type f -name "*.sh" -exec chown postgres:postgres {} +