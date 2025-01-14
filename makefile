build:
	docker build -t xdp-demo .

generate: build
	docker run --rm -v "$(PWD):/demo" --entrypoint go xdp-demo generate ./...

demo: build
	docker run --rm -it --name xdp-demo \
	  --cap-add SYS_RESOURCE --cap-add NET_ADMIN --cap-add BPF \
	  xdp-demo

shell: build
	docker run --rm -it --name xdp-demo --entrypoint sh \
	  --cap-add SYS_RESOURCE --cap-add NET_ADMIN --cap-add BPF \
	  xdp-demo
