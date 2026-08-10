.PHONY: up stop down clean \
	backend-up backend-stop \
	backend-deps backend-tidy backend-tools \
	backend-test backend-lint \
	frontend-up frontend-stop \
	admin-up admin-stop \
	preview-up preview-stop \
	logs ps

backend-lint:
	cd backend && golangci-lint run ./...

backend-deps:
	cd backend && go mod download

backend-tidy:
	cd backend && go mod tidy

backend-tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

backend-test:
	cd backend && go test -v -count=1 ./...

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
