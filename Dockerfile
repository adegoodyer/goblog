# First stage: Compile Go appllication
FROM golang:1.23 AS builder

# add labels to link container repo to code repo
LABEL org.opencontainers.image.source=https://github.com/adegoodyer/goblog
LABEL org.opencontainers.image.description="Simplest blog engine for coders"
LABEL org.opencontainers.image.licenses=MIT

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /goblog
COPY go.mod go.sum ./
RUN go mod download
COPY ./src ./src
COPY ./sample ./sample
RUN go build -o /goblog/goblog ./src

# Second stage: Create minimal image with compiled binary
FROM alpine:latest

WORKDIR /
COPY --from=builder /goblog/goblog /goblog
COPY --from=builder ./goblog/sample ./sample
ENTRYPOINT ["/goblog"]
