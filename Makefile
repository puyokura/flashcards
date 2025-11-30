.PHONY: build clean run

APP_NAME := flashcards
DIST_DIR := dist
SRC := .

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	@go build -o $(DIST_DIR)/$(APP_NAME) $(SRC)
	@cp habatan.csv $(DIST_DIR)/
	@cp DroidSansFallbackFull.ttf $(DIST_DIR)/
	@echo "Build complete. Artifacts in $(DIST_DIR)/"

clean:
	@echo "Cleaning up..."
	@rm -rf $(DIST_DIR)
	@echo "Clean complete."

run: build
	@echo "Running $(APP_NAME)..."
	@cd $(DIST_DIR) && ./$(APP_NAME)
