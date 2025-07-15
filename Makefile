
# this makefile is literally just a bash script. i just wanna get this shit done
# so here's the most boring and basic makefile you've ever seen

test: 
	@echo -e 'test failed'

build:
	cd ./src && go build -ldflags "-s -w" -trimpath -o mo main.go

install:
	cd ./src && sudo cp -v ./mo /usr/local/bin/mo

