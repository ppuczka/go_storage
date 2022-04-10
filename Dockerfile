# builder image
FROM golang:1.17 as builder
RUN mkdir /build
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -o go_storage 

# generate clean, final image for end users
FROM alpine:3.14


COPY --from=builder /src/go_storage .

# for container runing in kubernes additional files will be added from secrets
COPY --from=builder /src/config/*.yml ./config/
COPY --from=builder /src/tls/* ./tls/

EXPOSE 8088

# arguments that can be overridden
CMD [ "/go_storage" ]