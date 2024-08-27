# Load .env file
include ../.env

# Phony targets
.PHONY: mysql createdb dropdb createmigrate migrateup migratedown checkmigrate forcedowngrade

# Docker commands to manage MySQL container
mysql:
	docker run --name mysqlres -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) -d mysql

createdb:
	docker exec -it mysqlres mysql -u root -p$(MYSQL_ROOT_PASSWORD) -e "CREATE DATABASE $(DATABASE_NAME);"

dropdb:
	docker exec -it mysqlres mysql -u root -p$(MYSQL_ROOT_PASSWORD) -e "DROP DATABASE IF EXISTS $(DATABASE_NAME);"

# Migrate commands
createmigrate:
	migrate create -ext sql -dir ../migrations/mysql -seq init_schema

migrateup:
	@echo "Migrating up with database URL: $(MYSQL_URL)"
	@echo "DATABASE_NAME: $(DATABASE_NAME)"
	migrate -path ../migrations/mysql -database "$(MYSQL_URL)" -verbose up

migratedown:
	@echo "Migrating down with database URL: $(MYSQL_URL)"
	migrate -path ../migrations/mysql -database "$(MYSQL_URL)" -verbose down

checkmigrate:
	migrate -path ../migrations/mysql -database "$(MYSQL_URL)" version

forcedowngrade:
	@echo "Forcing downgrade to fix dirty state"
	migrate -path ../migrations/mysql -database "$(MYSQL_URL)" force 0
