
FROM golang:1.19
WORKDIR /go/src/github.com/cloudmarius/sample-go-serve-static
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM scratch
LABEL "org.opencontainers.image.source"="https://github.com/cloudmarius/sample-go-serve-static"
LABEL "org.opencontainers.image.description"="a sample app which serves and changes static files"
COPY --from=0 /go/src/github.com/cloudmarius/sample-go-serve-static/app /
EXPOSE 8080
ENTRYPOINT ["/app"]
