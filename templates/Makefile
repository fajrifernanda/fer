SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

testdbname=skeleton_test
test_command=richgo test ./... $(TEST_ARGS) -v --cover
migrate_up=go run main.go migrate --direction=up --step=0
migrate_down=go run main.go migrate --direction=down --step=0
run_command=go run main.go server
run_worker_command=go run main.go worker

proto:
	protoc --go_out=plugins=grpc:. pb/skeleton/*.proto
	ls pb/skeleton/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

repository/mock/mock_hello_repository.go:
	mockgen -destination=repository/mock/mock_hello_repository.go -package=mock gitlab.kumparan.com/yowez/skeleton-service/repository HelloRepository

client/mock/mock_client.go:
	mockgen -destination=client/mock/mock_client.go -package=mock gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton HelloServiceClient

worker/mock/mock_worker.go:
	mockgen -destination=worker/mock/mock_worker.go -package=mock gitlab.kumparan.com/yowez/skeleton-service/worker Worker

mockgen: repository/mock/mock_hello_repository.go worker/mock/mock_worker.go client/mock/mock_client.go

run: check-watchmedo-exists
	@watchmedo auto-restart -d . -R -p '*.go' -i '\.git' -- $(run_command)

run-worker: check-watchmedo-exists
	@watchmedo auto-restart -d . -R -p '*.go' -i '\.git' -- $(run_worker_command)

test:
	SVC_DISABLE_CACHING=true $(test_command)

check-watchmedo-exists:
	@watchmedo --version > /dev/null

migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run main.go migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi

clean:
	rm -v repository/mock/mock_*_repository.go
	rm -v mock/mock_*_client.go worker/mock/mock_*.go

.PHONY: run proto test clean mockgen check-watchmedo-exists 