FROM golang:1.19

# Объявляем аргументы
ARG JWTKEY=PXkGLtn0JZ5DFf4
ARG DBHOST=database
ARG DBPORT=5432
ARG DBNAME=site_database
ARG DBUSERNAME=postgres
ARG DBPASSWORD=postgres
ARG DBSSLMODE=disable
ARG TIMEZONE=Europe/Moscow


WORKDIR /app

COPY src/ .
RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]


