clear-postgres-data:
	docker compose down
	sudo rm -rf db/data

start-rcd:
	rclone rcd --rc-serve --rc-no-auth &
