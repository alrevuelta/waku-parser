# waku-parser

Prototype that parses the logs of a set of `nwaku` instances and a single `traffic` instance publishing traffic to estimate the propagation delay of each message. The tool works as follow:
* To be used with [nwaku-simulator](https://github.com/alrevuelta/nwaku-simulator), [waku-publisher](https://github.com/alrevuelta/waku-publisher) and [nwaku](https://github.com/waku-org/nwaku).
* This tool requires access to docker, that is used to fetch the log from all containers.
* A given `waku-publisher` instance (aka `traffic`) is used to inject the traffic into the network. Each message that is injected into the network is logged, with its hash and published timestamp. Right now the tool only allows one instance of the traffic injector.
* On the other hand, all `nwaku` containers are auto detected. Note that it must be compiled with a flag (TODO document) so that [this](https://github.com/waku-org/nwaku/blob/50412d18801ca4004894b87f7018cb21f4edeb14/waku/v2/node/waku_node.nim#L241) trace is printed upon rx a message.
* With all this information (logs from one `waku-publisher` and multiple `nwaku` instances) `waku-parser` parses the logs and prints every few seconds:
  * Amount of messages that were sent
  * Amount of messages that were received
  * Average delay for each interval
* Note that in order to mark a message as lost, a timeout can be configured.
* Metrics are also exposed as prometheus, see code. 

🚨Proof of concept, untested and dirty

Usage. Should work with default values, but if not, provide a `docker-host`:
```
$ go build
$ ./wakuparser --docker-host=unix:///Users/you-user/.docker/run/docker.sock --timeout-ms=5000
```