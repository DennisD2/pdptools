all:
	go build

clean:
	rm pdpcom

test:
	./pdpCom --batch --upload ./macro11-examples/054-multiecho.deposit

init:
	go mod init pdpcom
	go mod tidy