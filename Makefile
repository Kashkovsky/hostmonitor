monitor: clean-build client
	go build

prepare-release: clean client

client: clean-client
	cd web/client && yarn && yarn build && cp -R ./dist ../static && cd ../..

clean: clean-build clean-client

clean-build:
	rm -f hostmonitor
	rm -rf dist

clean-client:
	rm -rf web/static
	rm -rf web/client/dist
	rm -rf web/client/.parcel-cache
