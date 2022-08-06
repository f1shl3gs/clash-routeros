FROM golang:1.18.5-alpine3.16 as builder

RUN apk add --no-cache make git ca-certificates && \
    wget -O /Country.mmdb https://github.com/Dreamacro/maxmind-geoip/releases/latest/download/Country.mmdb
RUN wget -O /ui.zip https://github.com/Dreamacro/clash-dashboard/archive/refs/heads/gh-pages.zip && unzip /ui.zip -d /
WORKDIR /clash-src
COPY --from=tonistiigi/xx:golang / /
COPY . /clash-src

# Speed up mod download
ENV GOPROXY=https://goproxy.cn
RUN go mod download && \
    make build && \
    mv ./bin/clash /clash

FROM scratch

COPY --from=builder /Country.mmdb /
COPY --from=builder /clash /
COPY --from=builder /clash-dashboard-gh-pages /ui
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/clash"]
