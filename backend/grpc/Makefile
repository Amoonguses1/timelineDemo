GO_MODULE := grpc-test

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -fR ./protogen

.PHONY: protoc
protoc:
	protoc --go_opt=module=${GO_MODULE} --go_out=. \
		./proto/timeline/*.proto ./proto/post/*.proto

.PHONY: build
build: clean protoc tidy
