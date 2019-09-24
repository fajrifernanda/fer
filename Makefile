proto:
	@protoc --go_out=plugins=grpc:. pb/content/*.proto
	@ls pb/content/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

run:
	@go run main.go

build:
	@go build -o ./bin/trafo

changelog:
	@git-chglog -o CHANGELOG.md

