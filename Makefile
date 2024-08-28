build:
	@go build -o ./bin/fhir cmd/server/main.go

run: build
	@./bin/fhir