FROM golang:1.11 as build
WORKDIR /go/src/github.com/jukeizu/weather
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "jukeizu:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/jukeizu/weather/passwd /etc/passwd
COPY --from=build --chown=100:101 /go/src/github.com/jukeizu/weather/bin/weather .
USER jukeizu
ENTRYPOINT ["./weather"]
