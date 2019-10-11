FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/musclesister_api
COPY . .
RUN go build main.go

# runtime image
FROM alpine
COPY --from=builder /go/src/musclesister_api /app

CMD /app/main $PORT