.PHONY: build push

build:
	docker build --ssh default -t squizzi/msr-policy-updater .

push:
	docker push squizzi/msr-policy-updater:latest
