FROM mongo:6.0.13

COPY mongodb-keyfile /etc/mongo/mongodb-keyfile

RUN chmod 400 /etc/mongo/mongodb-keyfile

RUN chown 999:999 /etc/mongo/mongodb-keyfile

ENV MONGO_INITDB_ROOT_USERNAME=root

ENV MONGO_INITDB_ROOT_PASSWORD=password

CMD ["mongod", "--keyFile", "/etc/mongo/mongodb-keyfile", "--replSet", "rs0", "--bind_ip_all"]
