.PHONY: all clean build

define echo_green
    printf "\e[38;5;40m"
    echo "${1}"
    printf "\e[0m \n"
endef

current_dir=$(shell pwd)
goflags=GOFLAGS="-mod=readonly"
all: build

clean:
	rm -rf ./build

build: ${pb_go_files}
	${goflags} go build -o ./build/echo_template ./cmd/backend/...
	@${call echo_green,"build finished! The targets is in the ${current_dir}/build/"}

doc:
	swag init -d ./cmd/backend --parseDependency
	${goflags} go build -o ./build/doc ./cmd/api_doc/...
	@${call echo_green,"generate docs finished! The targets is in the ${current_dir}/build/"}