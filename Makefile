start.local:
	cd cmd/backend_coursework && go build -o app && ./app

up:
	cd deployments && docker compose up

up.detached:
	cd deployments && docker compose up -d

up.build:
	cd deployments && docker compose up --build

down:
	cd deployments && docker compose down
