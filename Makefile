docker:
	@docker buildx build -t sqlctest:latest .

run:
	docker compose down
	docker compose up -d
