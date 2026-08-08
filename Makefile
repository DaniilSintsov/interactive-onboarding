.PHONY: up stop down

up:
	docker compose up -d --build onboarding-backend

down:
	docker compose stop onboarding-backend db

clean:
	docker compose down
