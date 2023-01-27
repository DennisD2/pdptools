all:
	go build

clean:
	rm pdpcom

test:
	./pdpcom --batch --upload ../arthur-gill-examples/054-multiecho.deposit

init:
	go mod init pdpcom
	go mod tidy
