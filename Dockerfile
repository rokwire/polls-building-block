FROM golang:1.24-bullseye AS builder

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

RUN mkdir /polls-app
WORKDIR /polls-app
# Copy the source from the current directory to the Working Directory inside the container

COPY . .
RUN make

FROM alpine:3.22

#we need timezone database + certificates
RUN apk add --no-cache make tzdata ca-certificates


COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

COPY --from=builder /polls-app/bin/polls /
COPY --from=builder /polls-app/driver/web/docs/gen/def.yaml /driver/web/docs/gen/def.yaml

COPY --from=builder /polls-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /polls-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

COPY --from=builder /polls-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf /polls-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf
COPY --from=builder /polls-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf /polls-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf

ENTRYPOINT ["/polls"]