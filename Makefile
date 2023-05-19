CURRDIR:=$(shell pwd)
build-local:
	go build --buildmode=plugin -buildvcs=false -o lura-grpc-proxy.so ./cmd/lura-grpc-proxy
build-docker:
	docker run --rm -v $(CURRDIR):/build --workdir /build --entrypoint go devopsfaith/krakend-plugin-builder:2.1.2 build --buildmode=plugin -o lura-grpc-proxy.so  ./cmd/lura-grpc-proxy
build-testserver:
	docker build . -f build/test-server/Dockerfile -t testserver:latest
build-gateway:
	docker build . -f build/gateway/Dockerfile -t gateway:latest