MICROSERVICE_NAME=pablogod/counters
tag?=$(shell date +%Y%m%d-%H%M)
DOCKER_PUSH=docker --config ~/.docker/personal push

.PHONY: deploy run fmt vet

deploy:
	docker build -t $(MICROSERVICE_NAME):$(tag) -t $(MICROSERVICE_NAME):latest .
	$(DOCKER_PUSH) $(MICROSERVICE_NAME):$(tag)
	$(DOCKER_PUSH) $(MICROSERVICE_NAME):latest
	cd ~/work/personal/vps-deploy && make deploy-project p=counters
