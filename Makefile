build:
	@cd cmd/give-help-server; \
	go build -v

generate:
	rm -Rf generated
	mkdir generated
	swagger generate server -t generated  -P models.LoggedUser --exclude-main --skip-validation -f api/swagger.yml -r LICENSE

