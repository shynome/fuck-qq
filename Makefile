build:
	CGO_ENABLED=0 go build  -ldflags="-X 'main.Version=$$(git describe --tags --always --dirty)' -s -w" -o fuck-qq .
docker: build
	docker build . -t shynome/fuck-qq:$$(git describe --tags --always --dirty)
push: docker
	docker push shynome/fuck-qq:$$(git describe --tags --always --dirty)
