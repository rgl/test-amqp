dist: test-amqp.zip

build: test-amqp.exe

test-amqp.exe: *.go
	GOOS=windows GOARCH=amd64 go build -v -o $@ -ldflags="-s -w"

test-amqp.zip: test-amqp.exe
	rm -f $@
	zip $@ $^

clean:
	rm -rf test-amqp*

.PHONY: dist build clean
