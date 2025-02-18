.PHONY: build run migrate migrate-down seed clean reset

build:
	@mkdir -p bin
	go build -o bin/bills ./main.go
	go build -o bin/migrate ./cmd/migrate/main.go
	go build -o bin/seed ./cmd/seed/main.go

run: build
	./bin/bills

migrate: build
	./bin/migrate

migrate-down: build
	./bin/migrate -down

seed: build
	./bin/seed

clean:
	rm -f bills.db
	rm -rf bin/

# Full reset: clean DB and reseed
reset: clean build migrate seed 