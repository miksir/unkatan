FROM golang:1.14 AS build
ARG VER
ARG PRIVATE_USER
ARG PRIVATE_PASSWORD

COPY . /opt/

WORKDIR /opt/

#RUN echo -e "machine gitlab.private.com\nlogin ${PRIVATE_USER}\npassword ${PRIVATE_PASSWORD}" > ~/.netrc

RUN go build -o unkatan -ldflags "-X main.version=$VER" cmd/unkatan/main.go

FROM golang:1.14

COPY --from=build /opt/unkatan /opt/unkatan
COPY --from=build /opt/html /opt/html

WORKDIR /opt/

EXPOSE 8080

CMD ["/opt/unkatan", "--config", "/etc/config.yaml"]
