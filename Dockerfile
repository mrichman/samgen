FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/samgen

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/samgen .

FROM alpine:3.14
RUN apk add ca-certificates

COPY --from=build_base /tmp/samgen/out/samgen /app/samgen

CMD ["/app/samgen"]