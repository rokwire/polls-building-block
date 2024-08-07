FROM golang:1.22-bullseye as builder

ENV CGO_ENABLED=0

RUN mkdir /polls-app
WORKDIR /polls-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.20

#we need timezone database
RUN apk --no-cache add tzdata

COPY --from=builder /polls-app/bin/polls /
COPY --from=builder /polls-app/driver/web/docs/gen/def.yaml /driver/web/docs/gen/def.yaml

COPY --from=builder /polls-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /polls-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

#we need timezone database
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

ENTRYPOINT ["/polls"]
