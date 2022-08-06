VERSION   	:= $(shell grep "github.com/Dreamacro/clash" go.mod | awk '{print $$2}')
BUILDTIME	:= $(shell date -u)
GOBUILD		:= CGO_ENABLED=0 go build -tags timetzdata -trimpath -ldflags '-X "github.com/Dreamacro/clash/constant.Version=$(VERSION)" \
		-X "github.com/Dreamacro/clash/constant.BuildTime=$(BUILDTIME)" \
		-w -s -buildid='

# For local test
build:
	go mod download
	$(GOBUILD) -o bin/clash

docker:
	# make sure you have set up buildx correctly, see: https://docs.docker.com/build/buildx/multiplatform-images/
	docker buildx build -t clash:latest --platform linux/arm64 .