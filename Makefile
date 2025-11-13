migration:
ifndef name
	$(error name is not set. Usage: make migration name=[migration_name])
endif
	@num=$$(ls internal/db/migrations | wc -l | awk '{printf "%04d", $$1+1}'); \
	name="$(name)"; \
	touch internal/db/migrations/$${num}_$$name.sql; \
	echo "Created internal/db/migrations/$${num}_$$name.sql"