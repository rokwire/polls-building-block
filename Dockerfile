FROM golang:1.24-alpine as builder

ENV CGO_ENABLED=0

RUN apk add --no-cache tzdata make build-base

RUN mkdir /polls-app
WORKDIR /polls-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.22

#we need timezone database + certificates
RUN apk add --no-cache tzdata ca-certificates make

COPY --from=builder /polls-app/bin/polls /
COPY --from=builder /polls-app/driver/web/docs/gen/def.yaml /driver/web/docs/gen/def.yaml

COPY --from=builder /polls-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /polls-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

#we need timezone database
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

ENTRYPOINT ["/polls"]
