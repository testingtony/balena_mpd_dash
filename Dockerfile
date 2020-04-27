FROM golang:1.14.2-alpine as build

RUN apk update && apk add gcc libc-dev
WORKDIR /root
ADD . .
RUN go test ./...
RUN CGO_ENABLED=0 GOARCH=arm GOOS=linux go build -ldflags '-extldflags "-static"'

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch 

# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /root/dash_mpd .

ENTRYPOINT ["./dash_mpd"]