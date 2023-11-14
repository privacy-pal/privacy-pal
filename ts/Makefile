.PHONY: run_test, publish

run_test:
	tsc && node dist/test/chat/main.js

publish: run_test
	npm version $(keyword) --no-git-tag && npm publish
