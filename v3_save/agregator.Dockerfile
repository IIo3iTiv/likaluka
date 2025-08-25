FROM golang:1.25.0-alpine3.22 AS builder

WORKDIR /app
ENV USER=appuser
ENV UID=10001
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE="on"
ENV GOWORK="/app/go.work"
ENV GOPROXY="http://proxy.golang.org,direct"

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# RUN ls
COPY . .

RUN pwd
RUN go build -ldflags="-w -s" -o ./bin/main ./go-agregator/v2/

# ENTRYPOINT [ "/entypoint.sh" ]

FROM alpine:3.22.0 AS final

ENV USER=appuser
WORKDIR /app

COPY --chown=appuser --chmod=755 --from=builder /etc/passwd /etc/passwd
COPY --chown=appuser --chmod=755 --from=builder /etc/group /etc/group
COPY --chown=appuser --chmod=755 --from=builder /app/bin/main /app/bin/main
# COPY --chown=appuser --chmod=755 --from=builder /app/kafka/truststore/kafka.truststore.pem /app/bin/kafka.truststore.pem

EXPOSE 8080
ENTRYPOINT ["./bin/main"]