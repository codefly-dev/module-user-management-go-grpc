FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

COPY . .

COPY migrations /app/migrations

CMD /usr/local/bin/migrate -path /app/migrations -database "${CODEFLY__SERVICE_SECRET_CONFIGURATION__CUSTOMERS__STORE__POSTGRES__CONNECTION}" up
