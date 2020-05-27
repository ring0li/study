# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BUILDDIR=build

all: build-mac
# all: build-linux
	@cp app.yaml build/
	#@cp build/* ~/repos/
	#@cp -rfv hooks ~/repos/

# build-linux:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BUILDDIR)/gitflow_Linux app/main.go

build-mac:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -v -o $(BUILDDIR)/gitflow_Darwin app/main.go


