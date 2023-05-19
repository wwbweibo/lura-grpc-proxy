FROM devopsfaith/krakend-plugin-builder:2.1.2 AS build
WORKDIR /build
ADD . .
RUN go mod tidy
RUN go build -buildmode=plugin -o lura-grpc-proxy.so ./cmd/lura-grpc-proxy

FROM devopsfaith/krakend:2.1.2 AS publish
COPY --from=build /build/lura-grpc-proxy.so /opt/krakend/plugins/