# builder image
FROM golang:1.17 as builder
RUN mkdir /build
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -o go_storage 

# generate clean, final image for end users
FROM scratch
COPY --from=builder /src/go_storage .
COPY --from=builder /src/*.yml . 
COPY --from=builder /src/tls/ ./tls


EXPOSE 8088

# arguments that can be overridden
CMD [ "/go_storage" ]