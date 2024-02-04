Для запуска необходимо задать переменную окружения CONFIG_PATH, в которой бует путь к config.yaml.
Формат конфига:  
pgsql:   
  PG_USER: "postgres"  
  PG_PASSWORD: password  
  PG_DATABASE: "EWallet"  
  PG_PORT: "5432"  
  PG_HOST: "localhost"  
httpServer:  
  HTTP_HOST: "localhost"  
  HTTP_PORT: 80  
logger:  
  LOG_LEVEL: "local"  
  
Для запуска контейнера с postgresql можно воспользоваться:

FROM postgres:latest

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=mypassword
ENV POSTGRES_DB=EWallet

COPY /internal/repository/migrations/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["pgsql"]
