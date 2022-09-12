
FROM golang:1.19
WORKDIR /go/src/github.com/cloudmarius/sample-go-serve-static
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM scratch
LABEL key="value"
COPY --from=0 /go/src/github.com/cloudmarius/sample-go-serve-static/app /
EXPOSE 8080
ENTRYPOINT ["/app"]
