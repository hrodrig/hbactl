FROM golang:1.26-alpine AS build
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILDDATE=unknown
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/hrodrig/hbactl/cmd.Version=${VERSION}" -o /hbactl .

FROM alpine:3.23
LABEL org.opencontainers.image.title="hbactl"
LABEL org.opencontainers.image.description="CLI to manage PostgreSQL pg_hba.conf (list, add, remove, check, reload)"
LABEL org.opencontainers.image.source="https://github.com/hrodrig/hbactl"
LABEL org.opencontainers.image.authors="Hermes Rodríguez <https://github.com/hrodrig/hbactl>"
RUN apk --no-cache add ca-certificates
RUN adduser -D -g "" hbactl
COPY --from=build /hbactl /home/hbactl/hbactl
RUN chown hbactl:hbactl /home/hbactl/hbactl
USER hbactl
WORKDIR /home/hbactl
ENTRYPOINT ["/home/hbactl/hbactl"]
