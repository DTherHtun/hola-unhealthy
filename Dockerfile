FROM golang:alpine AS builder
RUN apk add --no-cache git && mkdir -p $GOPATH/src/github.com/DTherHtun/hola-unhealthy && go get github.com/rakyll/statik 
WORKDIR $GOPATH/src/github.com/DTherHtun/hola-unhealthy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/hola-unhealthy .
FROM scratch
COPY --from=builder /go/src/github.com/DTherHtun/hola-unhealthy/index.html /go/bin/index.html
COPY --from=builder /go/bin/hola-unhealthy /go/bin/hola-unhealthy
ENTRYPOINT ["/go/bin/hola-unhealthy"]
EXPOSE 8080
