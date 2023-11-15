proxy = proxy.golang.org
module = github.com/privacy-pal/privacy-pal/go

.PHONY: publish

publish:
	go test ./... && git tag $(version) && git push origin $(version) && GOPROXY=$(proxy) go list -m $(module)@${version}