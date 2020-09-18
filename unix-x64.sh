go build -a -installsuffix cgo -o bin/xftts main.go
export LD_LIBRARY_PATH=$(pwd)/xf/libs/x64