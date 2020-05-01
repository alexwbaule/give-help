build:
	go build -o bin/give-help-server -v cmd/give-help-server/main.go

generate:
	rm -Rf generated
	mkdir generated
	swagger generate server -t generated  -P models.LoggedUser --exclude-main --skip-validation -f api/swagger.yml -r LICENSE

importer-gdocs:
	go build -o bin/importer -v cmd/importer-gdocs/main.go

importer:
	go build -o bin/importer -v cmd/importer/main.go

exporter:
	go build -o bin/exporter -v cmd/exporter/main.go

all:
	rm -Rf generated
	mkdir generated
	swagger generate server -t generated  -P models.LoggedUser --exclude-main --skip-validation -f api/swagger.yml -r LICENSE
	go build -o bin/give-help-server -v cmd/give-help-server/main.go
	go build -o bin/exporter -v cmd/exporter/main.go
	go build -o bin/importer -v cmd/importer/main.go
	go build -o bin/importer-dgocs -v cmd/importer-gdocs/main.go
