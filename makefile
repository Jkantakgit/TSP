

.PHONY: all backend frontend clean



all: backend frontend

backend:
	@mkdir -p build
	@cd src/backend && go build -o ../../build/tsp-backend main.go
	@echo "✅ Built backend to build/tsp-backend"

frontend:
	@mkdir -p build
	@cd src/frontend && npm install
	@cd src/frontend && npx vite build
	@cp -r src/frontend/dist build/frontend
	@echo "✅ Built frontend to build/frontend"


clean:
	@rm -rf build
	@echo "🧹 Cleaned build directory"



dep:
	@cd src/backend && go mod tidy
	@cd src/frontend && npm install
	@echo "✅ Dependencies updated for backend and frontend"


run:
	@./build/tsp-backend & \
	backend_pid=$$!; \
	cd src/frontend && npm run dev; \
	kill $$backend_pid

