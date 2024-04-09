FROM golang:1.22 AS builder

RUN go install github.com/goreleaser/goreleaser@latest
WORKDIR /app
COPY ./.git/ ./.git/
RUN git reset --hard HEAD
RUN goreleaser build --clean --id=dns-grab --snapshot

FROM scratch
COPY --from=builder /app/dist/dns-grab_linux_amd64_v1/dns-grab /usr/local/bin/dns-grab
ENTRYPOINT [ "/usr/local/bin/dns-grab" ]