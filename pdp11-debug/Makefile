all:
	go build

clean:
	rm pdp11-debug

#test:
#	./pdp11-debug --batch --upload ../arthur-gill-examples/054-multiecho.deposit

init:
	go mod init pdp11-debug
	go mod tidy

install:
	go build
	cp -av pdp11-debug $(HOME)/bin