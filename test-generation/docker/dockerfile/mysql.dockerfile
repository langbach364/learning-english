FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=@ztegc4df9f4e
ENV MYSQL_DEFAULT_AUTHENTICATION_PLUGIN=mysql_native_password

COPY ./database/init.sql /docker-entrypoint-initdb.d/init.sql

CMD ["--init-file=/docker-entrypoint-initdb.d/init.sql"]