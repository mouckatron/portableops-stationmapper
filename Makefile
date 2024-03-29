
OL_VERSION=v6.14.1

build: build_ui embed_ui build_go

build_go:
	cd cmd/stationmapper && go build -o ../../build/stationmapper.exe
	cd cmd/wsjtx-mapper-ic && go build -o ../../build/wsjtx-mapper-ic.exe
	cd cmd/js8call-mapper-ic && go build -o ../../build/js8call-mapper-ic.exe

xcompile:
	cd cmd/stationmapper && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../../build/stationmapper.arm
	cd cmd/wsjtx-mapper-ic && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../../build/wsjtx-mapper-ic.arm
	cd cmd/js8call-mapper-ic && env GOOS=linux GOARCH=arm GOARM=6 go build -o ../../build/js8call-mapper-ic.arm

build_ui:
	cd ui && npm install
	cd ui/dist && rm -f *
	cd ui && npm run-script build

embed_ui:
	scripts/embed_ui.sh

clean:
	find . -type f -name '*~' -delete
