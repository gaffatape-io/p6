go_files :=  $(wildcard **/*.go)
go_test_files := $(wildcard **/*_test.go)
go_files += $(wildcard *.go)
go_non_test_files := $(filter-out $(go_test_files), $(go_files))

bin/p6: $(go_non_test_files)
	@echo $(go_non_test_files)
	go build -o bin/p6 ./srv/...


build: bin/p6


test: $(go_files)
	go test ./...


integration_test: $(go_files)
	env -u FIRESTORE_EMULATOR_HOST go test ./...


.PHONY:clean
clean:
	rm -rf bin/*
