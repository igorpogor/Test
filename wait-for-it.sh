#!/bin/bash
# wait-for-it.sh - script to wait for a service to be available

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

# Check if host and port are provided
if [ -z "$host" ] || [ -z "$port" ]; then
  echo "Usage: $0 host port -- command args"
  exit 1
fi

# Wait for the service to be available
until nc -w 1 "$host" "$port" >/dev/null 2>&1; do
  echo "Waiting for $host:$port..."
  sleep 1
done

# Execute the command
exec $cmd
