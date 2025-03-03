DEPLOYMENT_NAME := social-media-app
IMAGE_NAME := docker.io/library/social-media-app

TAG := $(shell date +"%Y%m%d-%H%M%S")

PREV_TAG := $(shell cat .last_image_tag 2>/dev/null)

build:
	echo "Building image: $(TAG)"
	kubectl scale deploy/$(DEPLOYMENT_NAME) --replicas=0
	sleep 3

	if [ ! -z "$(PREV_TAG)" ]; then \
		minikube image rm $(IMAGE_NAME):$(PREV_TAG) || true; \
	fi

	docker build -t $(IMAGE_NAME):$(TAG) .
	minikube image load $(IMAGE_NAME):$(TAG)

	kubectl set image deployment/$(DEPLOYMENT_NAME) $(DEPLOYMENT_NAME)=$(IMAGE_NAME):$(TAG)

	kubectl scale deploy/$(DEPLOYMENT_NAME) --replicas=1
	kubectl rollout restart deploy/$(DEPLOYMENT_NAME)

	echo "$(TAG)" > .last_image_tag
