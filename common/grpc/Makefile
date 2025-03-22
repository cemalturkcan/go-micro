PROTOC = protoc
PROTO_DIR = .

all: compile-protos

compile-protos:
	find $(PROTO_DIR) -name "*.proto" | xargs -I {} $(PROTOC) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. {}

clean:
	find $(PROTO_DIR) -name "*.pb.go" -delete