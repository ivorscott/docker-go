
api: 
	docker-compose up
	
cert:
	@echo "[ generating TLS certificate ]"
	go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	@mkdir -p tls;mv *.pem tls/