start-db:
	docker-compose up -d bot-db

stop-db:
	docker-compose down

build-linux:
	exec ./build.sh