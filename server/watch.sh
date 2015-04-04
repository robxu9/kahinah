#!/usr/bin/env bash

# usually invoked via fswatch
# fswatch -or -e '\.git/.*' -e '/build.*' -e '.*tmp.*' path/to/kahinah | xargs -n1 -I{} ./watch.sh

if [ -f /tmp/khserver.pid ]; then
    lastpid=$(cat /tmp/khserver.pid)
    if ps -p $lastpid > /dev/null; then
        kill -9 $lastpid # just kill it, we don't need to be nice
    fi
fi

touch /tmp/khserver.pid

./build.sh

(
    cd build
    ./khserver
    mv config.toml.new config.toml
    ./khserver &
    nextpid=$!
    echo "$nextpid" > /tmp/khserver.pid
)
