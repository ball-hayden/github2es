FROM golang:alpine AS builder

WORKDIR $GOPATH/src/haydenball.me.uk/github2es/

# Create appuser
ENV USER=appuser
ENV UID=1001

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' -a \
  -o /go/bin/github2es .

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/github2es /go/bin/github2es

USER ${USER}:${USER}

ENTRYPOINT ["/go/bin/github2es"]
