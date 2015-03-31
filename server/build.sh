#!/bin/bash

if ! hash ember 2>/dev/null; then
    echo "! requires ember-cli: install with npm install -g ember-cli"
fi

if ! hash go 2>/dev/null; then
    echo "! requires golang"
fi

echo "~> removing previous build directories"
rm -r build/

echo "~> clearing cache"
go clean -i -r -x

echo "~> creating build/"
mkdir -p build/

echo "~> generating server binary..."
go build -o build/khserver -i -race -v -x -a

echo "~> generating client files..."
(
  cd client
  ember build -prod --output-path ../build/public
)

echo "~> done."
