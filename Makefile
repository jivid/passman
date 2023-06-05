.PHONY: all cli

server:
	go build -o bin/server ./passman/

cli:
	go build -o bin/passman-cli ./cli/

clean:
	rm -f bin/*

all: clean cli server
