Curl in scratch image:
https://medium.com/axiomzenteam/combining-docker-multi-stage-builds-and-health-checks-feea7cd2d85e
https://github.com/tarampampam/curl-docker

https://github.com/MaksimDzhangirov/backendBankExample?tab=readme-ov-file

https://docs.docker.com/language/golang/build-images/

Миграции
https://copyprogramming.com/howto/how-to-run-golang-migrate-with-docker-compose#how-to-run-golang-migrate-with-docker-compose

https://github.com/peter-evans/docker-compose-actions-workflow?tab=readme-ov-file

https://peterevans.dev/posts/smoke-testing-containers/
https://github.com/mvdan/github-actions-golang
https://habr.com/ru/companies/otus/articles/650435/
https://habr.com/ru/articles/595627/
https://stackoverflow.com/questions/76488238/running-docker-and-github-actions-error-command-not-found
https://stackoverflow.com/questions/65883184/gokit-validate-request-payload-in-transport-layer

GRPC
https://github.com/travisjeffery/grpc-go-kit-error-example/blob/master/main.go
https://www.ru-rocker.com/2017/02/24/micro-services-using-go-kit-grpc-endpoint/
mustEmbedUnimplemented https://github.com/grpc/grpc-go/issues/3794
https://github.com/matryer/goblueprints/blob/master/chapter10/vault/server_grpc.go
https://habr.com/ru/articles/461279/
https://github.com/dayleader/golang-grpc-example/blob/master/cmd/main.go
https://habr.com/ru/articles/654645/

Project standard layout
https://github.com/golang-standards/project-layout

https://dev.to/nikl/how-to-build-a-containerized-microservice-in-golang-a-step-by-step-guide-with-example-use-case-5ea8
https://github.com/velotiotech/watermark-service/tree/master

Интересная статья про AGR и ENV в Dockerfile
https://vsupalov.com/docker-arg-env-variable-guide/

https://gokit.io/examples/stringsvc.html

Gokit: Validate request/payload in transport layer
https://stackoverflow.com/questions/65883184/gokit-validate-request-payload-in-transport-layer

Гайд про ENV
https://habr.com/ru/articles/446468/

Protobuf
1. Install (MacOS): brew install protobuf
2. Install: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
3. Install: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest