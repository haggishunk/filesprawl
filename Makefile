clear-postgres-data:
	docker container rm filesprawl-postgres-1
	docker volume rm filesprawl_postgres_data
