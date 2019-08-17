NAME := n
INSTALL := install
PKGROOT := github.com/phR0ze/$(NAME)

default: $(NAME)
$(NAME): $(INSTALL)
	go build -o bin/$(NAME) .

mech: $(DEPDIR)
	go build -o bin/mech pkg/net/mech/example/mech.go

install:
	dep ensure -v

update:
	dep ensure -v -update

test: $(NAME)
	@echo -e "\nRunning all go tests:"
	@echo -e "------------------------------------------------------------------------"
	go test $(PKGROOT)
	go test $(PKGROOT)/pkg/bin
	go test $(PKGROOT)/pkg/cli
	go test $(PKGROOT)/pkg/errs
	go test $(PKGROOT)/pkg/net
	go test $(PKGROOT)/pkg/opt
	go test $(PKGROOT)/pkg/sys
	go test $(PKGROOT)/pkg/tar
	go test $(PKGROOT)/pkg/term
	go test $(PKGROOT)/pkg/time
	go test $(PKGROOT)/pkg/tmpl
	go test $(PKGROOT)/pkg/unit

bench: $(NAME)
	@echo -e "\nRunning all go benchmarks:"
	@echo -e "------------------------------------------------------------------------"
	go test -bench=.

clean:
	rm -rf ./vendor
	rm -f ./bin/$(NAME)
