FROM --platform=$BUILDPLATFORM golang:1.27.3@sha256:f20d83a83e18dbd2a4714b8ffd4b30b287b6243962dd75608e86abcc63ba4efc AS build

ARG TARGETOS TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /glasfaser-modem-2-exporter ./cmd/glasfaser-modem-2-exporter/

FROM scratch

COPY --from=build /glasfaser-modem-2-exporter /glasfaser-modem-2-exporter

USER 65534:65534

EXPOSE 9100

ENTRYPOINT ["/glasfaser-modem-2-exporter"]
