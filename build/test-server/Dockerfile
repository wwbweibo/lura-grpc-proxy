FROM golang as build
WORKDIR /build
ADD . .
RUN go mod tidy
RUN go build -o test-server ./cmd/test-server

FROM ubuntu as publish
WORKDIR /app
EXPOSE 50051
COPY --from=build /build/test-server /app/
ENTRYPOINT [ "/app/test-server" ]