GO_MODULE := testGrpcTestClient/grpc

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -fR ./protogen

.PHONY: protoc
protoc:
	protoc --go_opt=module=${GO_MODULE} --go_out=. \
		--go-grpc_opt=module=${GO_MODULE} --go-grpc_out=. \
		./proto/timeline/*.proto ./proto/post/*.proto ./proto/image/*.proto

.PHONY: build
build: clean protoc tidy
