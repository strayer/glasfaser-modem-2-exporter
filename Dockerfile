FROM --platform=$BUILDPLATFORM golang:1.27.4@sha256:5f1b98fa02fe396ebce2a0227a0e08309e11792915b9dd2e1f71053d7e7e6666 AS build

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
