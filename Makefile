lint-breaking:
	buf breaking --against 'https://github.com/lucas-woo/cloud-drive'

lint-proto:
	buf lint --config buf.yaml

generate-proto:
	buf generate --template buf.gen.yaml

run-auth-server: cmd/grpc/auth/main.go;
	go run cmd/grpc/auth/main.go

run-iam-server: cmd/grpc/iam/main.go;
	go run cmd/grpc/iam/main.go

run-media-server: cmd/grpc/media/main.go;
	go run cmd/grpc/media/main.go

run-auth-test: cmd/test/auth/main.go
	go run cmd/test/auth/main.go;

run-iam-test: cmd/test/iam/main.go
	go run cmd/test/iam/main.go;

run-media-test: cmd/test/media/main.go
	go run cmd/test/media/main.go;