FROM --platform=$BUILDPLATFORM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS build

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
