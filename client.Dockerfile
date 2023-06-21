FROM golang:1.19

COPY . /app

WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /wow-client client/main.go

EXPOSE 8081

CMD ["/wow-client"]