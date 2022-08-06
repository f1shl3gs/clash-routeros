# clash-routeros
Clash container for RouterOS

This clash container is build for `RouterOS`, `UI` and `ExternalController` are enabled by
default, and listen on `0.0.0.0:9090`. Tested on Mikrotik RB5009.

## Build
[Buildx](https://docs.docker.com/build/buildx/multiplatform-images/) should be setup for build this image.
```shell
make docker
```

## Usage
1. Build container. Pre-build image is not provided, cause there is something need to be done.
```shell
# build docker
make docker

# save docker to local disk
docker save clash:latest > clash.tar
```

2. Upload the image to router.
  * Upload `clash.tar` with `WebFig`, aka http://192.168.88.1/webfig/#Files
  * Using ftp client
  * Copy `clash.tar` to USB stick, plug it to your router(if any)

3. Add necessary environments. Assume you have install container package and enable container mode,
see: https://help.mikrotik.com/docs/display/ROS/Container . 
```shell
/container/envs/add name=clash key=SUBSCRIPTION value="http://example.com"
/container/envs/add name=clash key=SUBSCRIPTION_UPDATE_INTERVAL value="6h"
```

4. Start clash container
```shell
/container/add file=clash.tar interface=veth1 envlist=clash hostname=clash logging=yes
# wait a while for extracting, it takes a few seconds, you can check the status of container by `/container print`

/container/start 0
```

5. Forward ports to internal Docker
```shell
# Forward proxy requests
/ip/firewall/nat/add action=dst-nat chain=dstnat dst-address=192.168.88.1 dst-port=7890 protocol=tcp to-addresses=172.17.0.2 to-ports=7890

# Forward controller requests
/ip/firewall/nat/add action=dst-nat chain=dstnat dst-address=192.168.88.1 dst-port=9090 protocol=tcp to-addresses=172.17.0.2 to-ports=9090
```

Now visit http://192.168.88.1/ui you shall see your clash dashboard

6. Setup your http proxy.
This is different from variant systems. If you known what clash is you should know how to do it. 
You may need to set environment `NO_PROXY=192.168.88.1` to avoid some proxy error.