#!/bin/sh

# Prerequisite: swagger is installed: github.com/go-swagger/go-swagger
curl -X GET "https://esi.evetech.net/latest/swagger.json" > swagger.json
swagger generate client --default-scheme=https --name=ESI --spec=swagger.json --copyright-file=copyright.txt --target=./
go get -u -f ./...
echo "Applying patch to mal-generated ./client/routes/get_route_origin_destination_parameters.go"
patch -p0 < get_route_origin_destination_parameters.patch
go build ./...
