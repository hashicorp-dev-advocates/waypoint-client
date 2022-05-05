PROTOC_ZIP := protoc-3.14.0-linux-x86_64.zip

setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Setup the dependencies for Mac
setup_mac: setup
	brew install protobuf

setup_linux: setup
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/${PROTOC_ZIP}
	sudo unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc
	sudo unzip -o ${PROTOC_ZIP} -d /usr/local 'include/*'
	rm -f ${PROTOC_ZIP}

generate_protos: get_google_protos get_waypoint_protos
	protoc \
		--go_out=./pkg \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pkg \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./protos \
		protos/waypoint/waypoint.proto

# Generate the mocks for the waypoint client using mockery
generate_mocks:
	docker run -it -v ${PWD}:/code -w /code/pkg/waypoint vektra/mockery

get_waypoint_protos:
	curl -L -s -o ./protos/waypoint/waypoint.proto https://raw.githubusercontent.com/hashicorp/waypoint/main/pkg/server/proto/server.proto
	curl -L -s -o ./protos/any.proto https://raw.githubusercontent.com/hashicorp/waypoint/main/thirdparty/proto/opaqueany/any.proto

get_google_protos:
	curl -L -s -o ./protos/google/protobuf/empty.proto https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/empty.proto
	curl -L -s -o ./protos/google/protobuf/timestamp.proto https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/timestamp.proto
	curl -L -s -o ./protos/google/protobuf/any.proto https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/any.proto
	curl -L -s -o ./protos/google/rpc/status.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/status.proto