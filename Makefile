DEPLOYMENT_NAME := social-media-app
IMAGE_NAME := social-media-app

MAJOR_VERSION := v1
MINOR_VERSION := 0
PATCH_VERSION := 0

TAG := $(MAJOR_VERSION).$(MINOR_VERSION).$(PATCH_VERSION)

build:
	echo "Building image: $(TAG)"
	docker build -t $(IMAGE_NAME):$(TAG) .
	minikube image load $(IMAGE_NAME):$(TAG)
	kubectl rollout restart deploy/$(DEPLOYMENT_NAME)
