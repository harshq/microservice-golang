install_swagger:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install_swagger
	swagger generate spec -o ./swagger.yaml --scan-models