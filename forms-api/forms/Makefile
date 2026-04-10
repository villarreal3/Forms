.PHONY: swagger-validate build

# Validar OpenAPI 2.0 con go-swagger (CLI local, Docker o Podman)
swagger-validate:
	@if command -v swagger >/dev/null 2>&1; then \
		swagger validate docs/swagger.yaml; \
	elif command -v docker >/dev/null 2>&1; then \
		docker run --rm -v "$$(pwd):/work" -w /work quay.io/goswagger/swagger validate docs/swagger.yaml; \
	elif command -v podman >/dev/null 2>&1; then \
		podman run --rm -v "$$(pwd):/work" -w /work quay.io/goswagger/swagger validate docs/swagger.yaml; \
	else \
		echo "Instale: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; \
		echo "O use Docker/Podman con imagen quay.io/goswagger/swagger"; \
		exit 1; \
	fi

build:
	go build -o app .
