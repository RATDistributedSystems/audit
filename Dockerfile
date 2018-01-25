FROM cassandra:latest

RUN mkdir /scripts
COPY /setup/create_database_structure.cql /setup/create_database_structure.sh /scripts/
WORKDIR "/scripts"
EXPOSE 9042
CMD ["./create_database_structure.sh"]