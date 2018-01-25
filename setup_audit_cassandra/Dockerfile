FROM cassandra:latest

RUN mkdir /scripts
COPY create_database_structure.cql create_database_structure.sh /scripts/
WORKDIR "/scripts"
EXPOSE 9042
CMD ["./create_database_structure.sh"]