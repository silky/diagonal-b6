all: protos experimental
	cd src/diagonal.works/diagonal; go build diagonal.works/diagonal/...
	cd src/diagonal.works/diagonal/cmd/fe; go build
	cd src/diagonal.works/diagonal/cmd/osm; go build
	cd src/diagonal.works/diagonal/cmd/osmbeam; go build
	cd src/diagonal.works/diagonal/cmd/inspect; go build
	cd src/diagonal.works/diagonal/cmd/splitosm; go build
	make -C data

docker: protos
	mkdir -p docker/bin/linux-amd64
	cd src/diagonal.works/diagonal/cmd/osm; GOOS=linux GOARCH=amd64 go build -o ../../../../../docker/bin/linux-amd64/osm
	cd src/diagonal.works/diagonal/cmd/splitosm; GOOS=linux GOARCH=amd64 go build -o ../../../../../docker/bin/linux-amd64/splitosm
	docker build -f docker/Dockerfile.diagonal -t diagonal docker
	docker tag diagonal eu.gcr.io/diagonal-platform/diagonal
	docker push eu.gcr.io/diagonal-platform/diagonal
	docker build -f docker/Dockerfile.monitoring -t monitoring docker
	docker tag monitoring eu.gcr.io/diagonal-platform/monitoring
	docker push eu.gcr.io/diagonal-platform/monitoring

protos:
	protoc -I=proto --go_out=src proto/geography.proto
	protoc -I=proto --go_out=src proto/tiles.proto
	protoc -I=proto --go_out=src proto/osm.proto
	protoc -I=src/diagonal.works/diagonal/osm --go_out=src src/diagonal.works/diagonal/osm/import.proto
	protoc -I=src/diagonal.works/diagonal/osm/pbf --go_out=src src/diagonal.works/diagonal/osm/pbf/pbf.proto

experimental:
	cd src/diagonal.works/diagonal/experimental/mr; go build
	cd src/diagonal.works/diagonal/experimental/osmpbf; go build

test:
	cd src/diagonal.works/diagonal; go test -v diagonal.works/diagonal/...

clean:
	find . -type f -perm +a+x | xargs rm

