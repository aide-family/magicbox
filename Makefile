.PHONY: init
init:
	@echo "Initializing magicbox environment"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.3
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/moon-monitor/stringer@latest
	go install github.com/protoc-gen/i18n-gen@latest


.PHONY: errors
# generate errors
errors:
	protoc --proto_path=./merr \
           --proto_path=./third_party \
           --go_out=paths=source_relative:./merr \
           --go-errors_out=paths=source_relative:./merr \
           ./merr/*.proto

.PHONY: generate
# generate stringer
generate:
	@echo "Generating stringer"
	go generate ./...
