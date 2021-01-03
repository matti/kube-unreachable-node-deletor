FROM golang:1.15.0-alpine3.12 as builder

COPY . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kube-unreachable-node-deletor .

FROM scratch
COPY --from=builder /app/kube-unreachable-node-deletor /kube-unreachable-node-deletor
ENTRYPOINT [ "/kube-unreachable-node-deletor", "5m"]
