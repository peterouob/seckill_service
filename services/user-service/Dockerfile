FROM golang:1.26 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

ARG TARGETARCH
RUN --mount=type=bind,source=.,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -v -o /main ./services/user-service/cmd/main.go

FROM gcr.io/distroless/static-debian13
WORKDIR /
COPY --from=base /main .

EXPOSE 8083 50050
ENTRYPOINT ["./main"]
