proxy = proxy.golang.org
module = github.com/privacy-pal/privacy-pal/go
subdir = go

.PHONY: publish, build_genpal, install_genpal

publish:
	go test ./... && git tag $(subdir)/$(version) && git push origin $(subdir)/$(version) && GOPROXY=$(proxy) go list -m $(module)@${version}

build_genpal:
	go build cmd/genpal/genpal.go

install_genpal:
	go install cmd/genpal/genpal.go