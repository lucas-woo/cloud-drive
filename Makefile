lint-breaking:
	buf breaking --against 'https://github.com/lucas-woo/cloud-drive'

lint-proto:
	buf lint --config buf.yaml

generate-proto:
	buf generate --template buf.gen.yaml

run-auth-server: cmd/grpc/auth/main.go;
	go run cmd/grpc/auth/main.go

run-auth-test: cmd/test/auth/main.go
	go run cmd/test/auth/main.go;