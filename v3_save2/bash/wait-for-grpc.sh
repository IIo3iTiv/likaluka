#!/bin/sh

host="$1"
port="$2"
shift 2
cmd="$@"

until grpc_health_probe -addr="$host:$port"; do
  >&2 echo "GRPC server $host:$port is unavailable - sleeping"
  sleep 2
done

>&2 echo "GRPC server $host:$port is up - executing command"
exec $cmd