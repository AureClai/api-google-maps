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
	rm ./statik -r -f
	rm ./bin/*.exe -r -f

clean:
	rm ./web/build -r -f
	rm ./statik -r -f
	rm ./bin/*.exe -r -f