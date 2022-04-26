
OL_VERSION=v6.14.1

build: deps embed_ui build_go

build_go:
	cd cmd/stationmapper && go build

deps:
	curl -sLC - --output lib/$(OL_VERSION)-dist.zip https://github.com/openlayers/openlayers/releases/download/$(OL_VERSION)/$(OL_VERSION)-dist.zip
	unzip -o lib/$(OL_VERSION)-dist.zip -d ui/
	mv ui/$(OL_VERSION)-dist/* ui/
	rmdir ui/$(OL_VERSION)-dist

embed_ui:
	scripts/embed_ui.sh

clean:
	find . -type f -name '*~' -delete
