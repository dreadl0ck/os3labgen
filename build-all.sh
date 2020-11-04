rm -rf bin/*
mkdir bin/darwin
mkdir bin/linux
mkdir bin/windows

GOOS=windows go build -o bin/windows/os3labgen.exe
GOOS=darwin go build -o bin/darwin/os3labgen
GOOS=linux go build -o bin/linux/os3labgen
