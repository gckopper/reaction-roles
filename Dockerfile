FROM golang:1.20-alpine AS build

WORKDIR /src

COPY *.go /src/
COPY go.mod /src/
COPY go.sum /src/

RUN go build

FROM scratch

COPY --from=build /src/reaction-roles /reaction-roles
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/reaction-roles"]
