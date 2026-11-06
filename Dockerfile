FROM --platform=$BUILDPLATFORM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e AS build

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
