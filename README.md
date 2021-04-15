# citus-failover
Worker failover support for citus community version using pg_auto_failover.

## What is this?

This is a simple service to monitor changes in a [pg_auto_failover][1] monitor and update workers' addresses in [citus][2] coordinator. The community version of citus does not support failover so you have to update the workers manually in case of a failover.

## How it works?

It will check both the citus coordinator and pg_auto_failover monitor and if it finds a worker node that is not primary, it will replace it with the primary one.

## Usage

You need to provide the connection details for both your coordinator and your monitor. Note that the coordinator should also be monitored by pg_auto_failover and you should provide the formation that is responsible for your coordinator, the primary node then will be found and used to update the workers.

You need to provide a config file. Check [`config.example.yml`][3] for a sample config.

**Monitor**:
PSQL connection details for your monitor.

**Coordinator**:
You need to provide the formation for your coordinator nodes as `Formation` instead of any host. It will be used to find the primary coordinator. 

### Run using binaries

```
./citus-failover -f path/to/config.yml
```

### Run using docker

Create a config file somewhere in your host. Then mount the path as `/config` to the container.

Map the port for API if enabled.

```
docker run -d -v /path/to/folder/for/config/file:/config -p 3002:3002 docker.pkg.github.com/navid2zp/citus-failover/citus-failover:latest
```

## REST API

A simple REST API is available. You can use it to check the state of your nodes. Set `API.Enabled` to `true` in your config file to enable the API. You also need to provide a secret string for the API. Then you'll need to provide this secret string in the header of your requests as `SECRET` to be authenticated and authorized.



Available endpoints:

`/v1/nodes`: List of all the nodes in your monitor.

`/v1/coordinator`: Primary coordinator address.

`/v1/coordinators`: List of all available coordinator nodes.

`v/1/workers/`: List of primary worker nodes.


[1]: https://github.com/citusdata/pg_auto_failover
[2]: https://github.com/citusdata/citus
[3]: https://github.com/Navid2zp/citus-failover/blob/main/config.example.yml
