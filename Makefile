.PHONY: build run clean

build:
	docker-compose up --build
	docker-compose run --rm migrate

run:
	docker-compose up --build -d

clean:
	docker-compose down