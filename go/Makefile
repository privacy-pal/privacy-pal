proxy = proxy.golang.org
module = github.com/privacy-pal/privacy-pal/go
subdir = go

.PHONY: publish

publish:
	go test ./... && git tag $(subdir)/$(version) && git push origin $(subdir)/$(version) && GOPROXY=$(proxy) go list -m $(module)@${version}