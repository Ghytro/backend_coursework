start.local:
	cd cmd/backend_coursework && go build -o app && cd ../../ && \
	DB_URL="postgres://postgres:mydbpassword@dockerdev.db:5432/postgres?sslmode=disable&" \
	./cmd/backend_coursework/app

start.db:
	cd deployments && docker compose up -d db

up:
	cd deployments && docker compose up

up.detached:
	cd deployments && docker compose up -d

up.build:
	cd deployments && docker compose up --build

down:
	cd deployments && docker compose down
