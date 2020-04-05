build:
	go build -o bin/give-help-server -v cmd/give-help-server/main.go

generate:
	rm -Rf generated
	mkdir generated
	swagger generate server -t generated  -P models.LoggedUser --exclude-main --skip-validation -f api/swagger.yml -r LICENSE

