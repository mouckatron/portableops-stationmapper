
OL_VERSION=v6.14.1

build: build_ui embed_ui build_go

build_go:
	cd cmd/stationmapper && go build -o ../../build/stationmapper.exe
	cd cmd/wsjtx-mapper-ic && go build -o ../../build/wsjtx-mapper-ic.exe

xcompile:
	cd cmd/stationmapper && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../../build/stationmapper.arm
	cd cmd/wsjtx-mapper-ic && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../../build/wsjtx-mapper-ic.arm

deps:
	curl -sLC - --output lib/$(OL_VERSION)-dist.zip https://github.com/openlayers/openlayers/releases/download/$(OL_VERSION)/$(OL_VERSION)-dist.zip
	unzip -o lib/$(OL_VERSION)-dist.zip -d ui/dist/
	mv ui/dist/$(OL_VERSION)-dist/* ui/dist
	rmdir ui/dist/$(OL_VERSION)-dist

build_ui:
	cp ui/src/*.html ui/dist/
	cp ui/src/*.js ui/dist/

embed_ui:
	scripts/embed_ui.sh

clean:
	find . -type f -name '*~' -delete
