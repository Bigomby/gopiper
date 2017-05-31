MKL_RED?=	\033[031m
MKL_GREEN?=	\033[032m
MKL_YELLOW?=	\033[033m
MKL_BLUE?=	\033[034m
MKL_CLR_RESET?=	\033[0m

BIN=      gopiper
CMD=      ./cmd/$(BIN)
bindir?=	$(prefix)/bin

build: vendor
	@printf "$(MKL_YELLOW)[BUILD]$(MKL_CLR_RESET)    Building project\n"
	@go build -ldflags "-X main.version=`git describe --tags --always --dirty=-dev`" -o $(BIN) $(CMD)

install: build
	@printf "$(MKL_YELLOW)[INSTALL]$(MKL_CLR_RESET)  Installing $(BIN)\n"
	@go install -ldflags "-X main.version=`git describe --tags --always --dirty=-dev`" $(CMD)

uninstall:
	@printf "$(MKL_RED)[UNINSTALL]$(MKL_CLR_RESET)  Remove $(BIN) from $(bindir)\n"
	@rm $(bindir)/$(BIN)

PACKAGE_LIST := $(shell glide novendor)
tests: vendor
	@printf "$(MKL_GREEN)[TESTING]$(MKL_CLR_RESET)  Running tests\n"
	@go test -race -v $(PACKAGE_LIST)

coverage:
	@printf "$(MKL_BLUE)[COVERAGE]$(MKL_CLR_RESET) Computing coverage\n"
	@echo "mode: count" > coverage.out
	@go test -covermode=count -coverprofile=gopiper.cover ./
	@go test -covermode=count -coverprofile=component.cover ./component
	@grep -h -v "mode: count" *.cover >> coverage.out
	@go tool cover -func coverage.out

GLIDE := $(shell command -v glide 2> /dev/null)
vendor:
ifndef GLIDE
	$(error glide is not installed)
endif
	@printf "$(MKL_BLUE)[DEPS]$(MKL_CLR_RESET)  Resolving dependencies\n"
	@glide install

clean:
	rm -f $(BIN)
	rm -f coverage.out

vendor-clean:
	rm -rf vendor/

all: build tests coverage
