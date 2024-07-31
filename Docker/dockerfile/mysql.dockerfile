FROM mysql:latest

COPY ../Handler/Database/init.sql /docker-entrypoint-initdb.d/init.sql
