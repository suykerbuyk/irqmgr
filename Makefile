
all: clean build

build: build_client build_server

build_server:
	go build -o bin/irqmgr_server server/irqmgrServer.go
build_client:
	go build -o bin/irqmgr_client client/irqMgrClient.go client/jsonui.go client/tree.go

clean:
	rm bin/*
