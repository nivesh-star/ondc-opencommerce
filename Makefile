# Default target to build the Docker image
request-action-service:
	docker build --rm --no-cache -t request-action-service -f buyer-platform/request-action-service/Dockerfile .
buyer-app-service:
	docker build --rm --no-cache -t buyer-app-service -f buyer-platform/buyer-app-service/Dockerfile .

# # Target to remove the image
# clean:
# 	docker rmi $(IMAGE_NAME)

# # Target to prune dangling images
# prune:
# 	docker image prune -f

# # Target to view the images (optional)
# images:
# 	docker images