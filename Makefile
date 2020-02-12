.ONESHELL:
.PHONY: all

all: clean yarn statik build run

gobuild: cleango statik build run
	
yarn:
	cd web && yarn build && cd ..

statik:
	statik -src=./web/build

build:
	go build -o ./bin/testStatik.exe

run:
	./bin/testStatik.exe

cleango:
	rm -r -f ./statik
	rm -r -f ./bin/*.exe

clean:
	rm -r -f ./web/build 
	rm -r -f ./statik 
	rm -r -f ./bin/*.exe 