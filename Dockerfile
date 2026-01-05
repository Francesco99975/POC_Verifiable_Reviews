FROM golang:1.25.5-alpine3.21 AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/reviews ./cmd/server/*.go

FROM alpine:3.19 AS release

RUN apk update && apk upgrade && apk --no-cache add ca-certificates

WORKDIR /go/bin

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/static /go/bin/static
COPY --from=build /go/src/app/sql /go/bin/sql

EXPOSE 3169

ENTRYPOINT ["/go/bin/reviews", "--port", "3169"]
