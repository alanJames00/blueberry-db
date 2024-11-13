# variables
TESTS_DIR=./tests
BUILD_DIR=build
SRC_DIR=./cmd/server
BINARY_NAME=blueberrydb
PID_FILE=server.pid

# Composite build, run and test
build-run-test: build run-background test stop

# Build the project
build:
	@echo "Building the project"
	mkdir ${BUILD_DIR}
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

# Run the project
run: build
	@echo "Running the project"
	./$(BUILD_DIR)/$(BINARY_NAME)

# Test the project
test:
	@echo "Running Integration Tests"
	go test -v $(TESTS_DIR)

# Run background for Running tests in parallel
run-background:
	@echo "Running project in background"
	./$(BUILD_DIR)/$(BINARY_NAME) & echo $$! > $(PID_FILE)

# Stop the running project with pid file
stop:
	@echo "Stopping running project"
	@if [ -f $(PID_FILE) ]; then \
		kill `cat $(PID_FILE)`; \
		rm -f $(PID_FILE); \
		echo "Project stopped"; \
	else \
		echo "No running project found"; \
	fi
