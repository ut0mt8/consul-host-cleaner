FROM golang:1.21 AS builder

ARG GITHUB_TOKEN

WORKDIR /app
COPY . .

RUN make staticbuild

FROM alpine
COPY --from=builder /app/cleaner /app/cleaner

ENTRYPOINT ["/app/cleaner"]
