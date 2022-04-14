monitor: clean client
	go build
client:
	rm -rf web/static && cd web/client && rm -rf dist && yarn && yarn build && cp -R ./dist ../static && cd ../..
clean:
	rm -f hostmonitor
