db-up:
	docker compose up -d 

db-clean:
	docker compose down
	sudo rm -rf db/data

db-reset: db-clean db-up

start-rcd:
	rclone rcd --rc-serve --rc-no-auth &

build:
	go build -o build/filesprawl
