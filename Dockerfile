FROM --platform=$BUILDPLATFORM golang:1.27.4@sha256:5e8afd493a7124a8586e134d9b532731161d4b9becf9dee5a5956b107426f93d AS build

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
