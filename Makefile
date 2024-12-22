# Default target to build the Docker image
request-action-service:
	docker build --rm --no-cache -t request-action-service -f buyer-platform/request-action-service/Dockerfile .
buyer-app-service:
	docker build --rm --no-cache -t buyer-app-service -f buyer-platform/buyer-app-service/Dockerfile .
bap-api:
	docker build --rm --no-cache -t bap-api -f buyer-platform/bap-api/Dockerfile .
bap-adapter-service:
	docker build --rm --no-cache -t bap-adapter-service -f buyer-platform/bap-adapter-service/Dockerfile .

all: request-action-service buyer-app-service bap-api bap-adapter-service

# # Target to remove the image
# clean:
# 	docker rmi $(IMAGE_NAME)

# # Target to prune dangling images
# prune:
# 	docker image prune -f

# # Target to view the images (optional)
# images:
# 	docker images