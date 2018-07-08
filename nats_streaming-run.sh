#!/bin/sh
docker volume create natsstreamdata
docker run -it -h nats-streaming \
	--restart always \
	--name nats-streaming \
	-v natsstreamdata:/data \
	-p 4222:4222 \
	nats-streaming -st FILE --dir /data -mc 0 -mm 0 -mb 0
