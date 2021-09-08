FROM golang:1.17.0-alpine3.14 AS build

# alpine comes without gcc
ENV CGO_ENABLED=0

RUN mkdir /opt/servicedef-src
WORKDIR /opt/servicedef-src
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go test ./...
RUN go build -v

FROM golang:1.17.0-alpine3.14 AS runtime

ARG BUILD_DATE
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE

RUN apk add --update-cache nmap nmap-scripts
RUN mkdir /opt/servicedef
WORKDIR /opt/servicedef
COPY --from=build /opt/servicedef-src/servicedef /opt/servicedef/servicedef
COPY hosts.json.example services.json.example rules.json.example /opt/servicedef/

ENTRYPOINT ["/opt/servicedef/servicedef"]
