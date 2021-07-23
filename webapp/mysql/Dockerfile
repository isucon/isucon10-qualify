FROM mysql:5.7

COPY db/* /docker-entrypoint-initdb.d/

CMD ["--character-set-server=utf8mb4"]