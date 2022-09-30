# builder stage
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /App

COPY ./app/ ./app/
COPY ./go.mod .
COPY ./go.sum .
COPY ./main.go .
COPY ./config.yml .
# add more package/file as needed

RUN go mod download
RUN go build -o application main.go


#Package stage
FROM alpine:3.16

WORKDIR /App

COPY --from=builder /App/application .
COPY --from=builder /App/config.yml .

CMD ["/App/application -config=config.yml"]
