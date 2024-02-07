clear-postgres-data:
	docker container rm filesprawl-postgres-1
	docker volume rm filesprawl_postgres_data

start-rcd:
	rclone rcd --rc-serve --rc-no-auth &
