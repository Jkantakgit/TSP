.PHONY: all backend frontend clean dep run-backend run-frontend

BACKEND_DIR = src/backend
FRONTEND_DIR = src/frontend
BUILD_DIR = build
BACKEND_BIN = $(BUILD_DIR)/tsp-backend
FRONTEND_BUILD = $(BUILD_DIR)/frontend

all: backend frontend

backend:
	@mkdir -p $(BUILD_DIR)
	@cd $(BACKEND_DIR) && go build -o ../../$(BACKEND_BIN) main.go
	@echo "Built backend to $(BACKEND_BIN)"

frontend:
	@mkdir -p $(BUILD_DIR)
	@cd $(FRONTEND_DIR) && npm install
	@cd $(FRONTEND_DIR) && npx vite build
	@cp -r $(FRONTEND_DIR)/dist $(FRONTEND_BUILD)
	@echo "Built frontend to $(FRONTEND_BUILD)"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build directory"

dep:
	@cd $(BACKEND_DIR) && go mod tidy
	@cd $(FRONTEND_DIR) && npm install
	@echo "Dependencies updated for backend and frontend"

# Run backend separately
run-backend:
	@echo "Starting backend server..."
	@$(BACKEND_BIN)

# Run frontend dev server separately
run-frontend:
	@cd $(FRONTEND_DIR) && npm run dev
