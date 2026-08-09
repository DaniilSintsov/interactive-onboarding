.PHONY: up stop down clean clean-all \
	backend-up backend-stop \
	frontend-up frontend-stop \
	admin-up admin-stop \
	preview-up preview-stop \
	logs ps

lint:
	cd backend && golangci-lint run ./...

up:
	docker compose up -d --build

stop:
	docker compose stop

down:
	docker compose down --remove-orphans

clean:
	docker compose down --volumes --remove-orphans

backend-up:
	docker compose up -d --build onboarding-backend

backend-stop:
	docker compose stop onboarding-backend db

frontend-up:
	docker compose up -d --build admin test-preview

frontend-stop:
	docker compose stop admin test-preview

admin-up:
	docker compose up -d --build admin

admin-stop:
	docker compose stop admin

preview-up:
	docker compose up -d --build test-preview

preview-stop:
	docker compose stop test-preview

logs:
	docker compose logs -f --tail=100

ps:
	docker compose ps
