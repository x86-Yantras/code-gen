CONTAINER_IMAGE_NAME=code-gen-container-image

container-image: .container-image

.container-image: Containerfile go.mod go.sum $(shell find . -name '*.go') $(shell find templates)
	podman build --file $< --tag ${CONTAINER_IMAGE_NAME}
	podman image inspect --format "{{.Id}}" ${CONTAINER_IMAGE_NAME} > $@

container: container-image
	podman run --rm -it ${CONTAINER_IMAGE_NAME} bash
