FROM boxfuse/flyway:5

RUN apt-get update \
  && apt-get upgrade -y \
  && rm -rf /var/lib/apt/lists/*

ENV DATABASE_SCHEMAS cocoapods_stats_rolluped

RUN rm /flyway/sql/put-your-sql-migrations-here.txt
COPY sql/*.sql /flyway/sql/

COPY sql/run-migrations /usr/local/bin/run-migrations
RUN chmod +x /usr/local/bin/run-migrations

ENTRYPOINT ["/usr/local/bin/run-migrations"]
