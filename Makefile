hello:
	@echo "Hello"

build:
	docker buildx build -t localhost/sqlctest:latest .

run:
	docker run --rm -it -p 8080:8080 localhost/sqlctest:latest

all: hello build run