lint-breaking:
	buf breaking --against 'https://github.com/lucas-woo/cloud-drive'

lint-proto:
	buf lint --config buf.yaml

generate-proto:
	buf generate --template buf.gen.yaml
