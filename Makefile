GRPC_FILES := grpc/auth/authpb.proto grpc/biz/sqlpb.proto
SWAGGER_JSON_FILES = ./authpb.swagger.json ./sqlpb.swagger.json

generate-proto:
	mkdir -p gen
	$(foreach file, $(GRPC_FILES), \
		protoc --experimental_allow_proto3_optional \
		--go_out=./gen --go_opt=paths=source_relative \
  		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
  		$(file) && \
  		protoc --experimental_allow_proto3_optional \-I . --grpc-gateway_out ./gen \
  		--grpc-gateway_opt logtostderr=true \
  		--grpc-gateway_opt paths=source_relative \
  		--grpc-gateway_opt generate_unbound_methods=true \
  		--grpc-gateway_opt generate_unbound_methods=true \
  		$(file); \
  	)


generate-google-annotations:
	mkdir -p google/api
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto >	google/api/annotations.proto
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto >google/api/http.proto

generate-swagger-json:
	# Since protoc-gen-swagger does not support optional fields, remove the optional keyword in .proto file when running this command.
	$(foreach file, $(GRPC_FILES), \
    	protoc -I=. --swagger_out=logtostderr=true:. $(file); \
    )

generate-swagger-server:
	docker pull quay.io/goswagger/swagger
	alias goswagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
	goswagger generate server -f ./authpb.swagger.json

	$(foreach s, $(SWAGGER_JSON_FILES), \
    	goswagger generate server -f $(swagger_file); \
    )

clean-generated:
	rm gen/grpc/*.go

swag-init:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init
