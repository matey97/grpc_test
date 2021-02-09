FROM golang:1.15-alpine as builder

WORKDIR /app

RUN mkdir server && mkdir grpc_test

COPY ./grpc_test/go.* ./grpc_test/
RUN cd grpc_test && go mod download && cd ..

COPY ./server/go.* ./server/
RUN cd server && go mod download && cd ..

COPY ./grpc_test/* ./grpc_test/
COPY ./server/* ./server/

RUN cd server && go build -v -o server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server/server /server
COPY google_services.json /tmp/keys/google_services.json

EXPOSE 443
CMD [ "/server" ]