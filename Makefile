.PHONY: container
container:
	# http://www.blang.io/2015/04/19/golang-alpine-build-golang-binaries-for-alpine-linux.html
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo multiset.go
	docker build -t elsdoerfer/consul2vulcand .
