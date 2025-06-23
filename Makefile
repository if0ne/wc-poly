PROJECT_NAME = exp
COMPONENTS = cmp prv

.PHONY: all clean build deploy

all: build deploy

clean:
	wash app delete wadm.yaml

build:
	@echo "Starting build process..."
	@for component in $(COMPONENTS); do \
		echo "Building $$component..."; \
		wash build -p $$component; \
	done
	@echo "All components built successfully."

deploy:
	@echo "Deploying..."
	wash app deploy wadm.yaml
	@echo "Deployment completed!"
