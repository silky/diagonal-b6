# Sets TARGETARCH to something linke amd64 or aarch64
TARGETARCH ?= $(shell uname -m | tr A-Z a-z)
# Sets TARGETOS to something like linux or darwin
TARGETOS ?= $(shell uname -s | tr A-Z a-z)
# Sets TARGETPLATFORM to something like linux/amd64 or darwin/aarch64
export TARGETPLATFORM ?= ${TARGETOS}/${TARGETARCH}

all: protos experimental fe ingest transit fe-js dfe
	cd src/diagonal.works/diagonal/monitoring; go generate
	cd src/diagonal.works/diagonal; go build diagonal.works/diagonal/...
	cd src/diagonal.works/diagonal/cmd/inspect; go build
	cd src/diagonal.works/diagonal/cmd/splitosm; go build
	cd src/diagonal.works/diagonal/cmd/tile; go build
	cd src/diagonal.works/diagonal/experimental/sightline-tiles; go build
	make -C data

fe: protos
	cd src/diagonal.works/diagonal/monitoring; go generate
	cd src/diagonal.works/diagonal/cmd/fe; go build -o ../../../../../bin/${TARGETPLATFORM}/fe

fe-js:
	make -C js

ingest: protos
	cd src/diagonal.works/diagonal/monitoring; go generate
	cd src/diagonal.works/diagonal/cmd/ingest; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest

ingest-gdal:
	cd src/diagonal.works/b6/cmd/ingest-gdal; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest-gdal

ingest-terrain:
	cd src/diagonal.works/b6/cmd/ingest-terrain; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest-terrain

ingest-gb-uprn:
	cd src/diagonal.works/b6/cmd/ingest-gb-uprn; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest-gb-uprn

ingest-gb-codepoint:
	cd src/diagonal.works/b6/cmd/ingest-gb-codepoint; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest-gb-codepoint

ingest-gb-ons:
	cd src/diagonal.works/b6/cmd/ingest-gb-ons; go build -o ../../../../../bin/${TARGETPLATFORM}/ingest-gb-ons

connect:
	cd src/diagonal.works/b6/cmd/connect; go build -o ../../../../../bin/${TARGETPLATFORM}/connect

b6-api:
	cd src/diagonal.works/b6/cmd/b6-api; go build -o ../../../../../bin/${TARGETPLATFORM}/b6-api

transit: protos
	cd src/diagonal.works/diagonal/cmd/transit; go build

mbtiles:
	cd src/diagonal.works/diagonal/cmd/mbtiles; go build

tile-profile:
	cd src/diagonal.works/diagonal/cmd/tile-profile; go build -o ../../../../../bin/${TARGETPLATFORM}/tile-profile

baseline: protos src/diagonal.works/b6/api/y.go	
	make -C src/diagonal.works/diagonal/cmd/baseline

baseline-backend: protos src/diagonal.works/b6/api/y.go
	make -C src/diagonal.works/diagonal/cmd/baseline baseline

b6: protos src/diagonal.works/b6/api/y.go
	cd src/diagonal.works/b6/cmd/b6; go build -o ../../../../../bin/${TARGETPLATFORM}/b6
	make -C src/diagonal.works/b6/cmd/b6/js

dfe:
	mkdir -p bin/${TARGETPLATFORM}
	cd src/diagonal.works/diagonal/cmd/dfe; go build -o ../../../../../bin/${TARGETPLATFORM}/dfe

tiles: protos
	mkdir -p bin/${TARGETPLATFORM}
	cd src/diagonal.works/diagonal/cmd/tiles; go build -o ../../../../../bin/${TARGETPLATFORM}/tiles

docker: protos
	mkdir -p docker/bin/linux-amd64
	cd src/diagonal.works/diagonal/cmd/ingest; GOOS=linux GOARCH=amd64 go build -o ../../../../../docker/bin/linux-amd64/ingest
	cd src/diagonal.works/diagonal/cmd/splitosm; GOOS=linux GOARCH=amd64 go build -o ../../../../../docker/bin/linux-amd64/splitosm
	docker build -f docker/Dockerfile.diagonal -t diagonal docker
	docker tag diagonal eu.gcr.io/diagonal-platform/diagonal
	docker push eu.gcr.io/diagonal-platform/diagonal
	docker build -f docker/Dockerfile.monitoring -t monitoring docker
	docker tag monitoring eu.gcr.io/diagonal-platform/monitoring
	docker push eu.gcr.io/diagonal-platform/monitoring
	docker build -f docker/Dockerfile.planet -t planet docker
	docker tag planet eu.gcr.io/diagonal-platform/planet
	docker push eu.gcr.io/diagonal-platform/planet

docker-build:
	docker build -f docker/Dockerfile.build -t build docker

docker-atlas-dev-data:
	cp data/earth/ne_10m_land.shp docker/data/atlas-dev
	cp data/earth/ne_10m_land.prj docker/data/atlas-dev
	docker build -f docker/Dockerfile.atlas-dev-data -t atlas-dev-data docker
	docker tag atlas-dev-data eu.gcr.io/diagonal-platform/atlas-dev-data
	docker push eu.gcr.io/diagonal-platform/atlas-dev-data

docker-ingest:
	scripts/make-in-docker.sh ingest
	mkdir -p docker/bin/${TARGETPLATFORM}
	cp bin/${TARGETPLATFORM}/ingest docker/bin/${TARGETPLATFORM}
	docker build --build-arg platform=${TARGETPLATFORM} -f docker/Dockerfile.ingest -t ingest-${TARGETARCH} docker
	docker tag ingest-${TARGETARCH} eu.gcr.io/diagonal-platform/ingest-${TARGETARCH}
	docker push eu.gcr.io/diagonal-platform/ingest-${TARGETARCH}

# Use TARGETARCH=x86_64 TARGETOS=linux for GCP
docker-atlas-dev: fe-js docker-atlas-dev-data
	mkdir -p docker/bin/${TARGETPLATFORM}
	scripts/make-in-docker.sh fe
	cp bin/${TARGETPLATFORM}/fe docker/bin/${TARGETPLATFORM}
	mkdir -p docker/js
	rm -rf docker/js/dist
	cp -r js/dist docker/js/dist
	docker build -f docker/Dockerfile.atlas-dev -t atlas-dev docker
	docker tag atlas-dev eu.gcr.io/diagonal-platform/atlas-dev
	docker push eu.gcr.io/diagonal-platform/atlas-dev

docker-dfe:
	mkdir -p docker/bin/${TARGETPLATFORM}
	cp bin/${TARGETPLATFORM}/dfe docker/bin/${TARGETPLATFORM}
	docker build --build-arg platform=${TARGETPLATFORM} -f docker/Dockerfile.dfe -t dfe-${TARGETARCH} docker
	docker tag dfe-${TARGETARCH} eu.gcr.io/diagonal-platform/dfe-${TARGETARCH}
	docker push eu.gcr.io/diagonal-platform/dfe-${TARGETARCH}

docker-tiles:
	mkdir -p docker/bin/${TARGETPLATFORM}
	cp bin/${TARGETPLATFORM}/tiles docker/bin/${TARGETPLATFORM}
	docker build --build-arg platform=${TARGETPLATFORM} -f docker/Dockerfile.tiles -t tiles-${TARGETARCH} docker
	docker tag tiles-${TARGETARCH} eu.gcr.io/diagonal-platform/tiles-${TARGETARCH}
	docker push eu.gcr.io/diagonal-platform/tiles-${TARGETARCH}

data/region/scottish-borders.index:
	gsutil cp gs://diagonal.works/region/scottish-borders.index data/region/scottish-borders.index

data/region/scottish-borders.connected.overlay:
	gsutil cp gs://diagonal.works/region/scottish-borders.connected.overlay data/region/scottish-borders.connected.overlay

# Use TARGETARCH=x86_64 TARGETOS=linux for GCP
docker-baseline:
	rm -rf docker/baseline
	mkdir -p docker/bin/${TARGETPLATFORM}
	cp bin/${TARGETPLATFORM}/baseline docker/bin/${TARGETPLATFORM}
	mkdir -p docker/baseline/assets/fonts
	cp -r js/dist/fonts/national-* docker/baseline/assets/fonts
	cp -r js/dist/fonts/unica77-* docker/baseline/assets/fonts
	mkdir -p docker/baseline/assets/images
	cp js/dist/images/logo.svg docker/baseline/assets/images
	cp js/dist/images/zoom-in.svg docker/baseline/assets/images
	cp js/dist/images/zoom-out.svg docker/baseline/assets/images
	cp js/dist/images/zoom-out.svg docker/baseline/assets/images
	mkdir -p docker/baseline/assets/icons
	cp js/dist/icons/favicon-16x16.png docker/baseline/assets/icons
	cp js/dist/icons/favicon-32x32.png docker/baseline/assets/icons
	cp js/dist/icons/hardware.svg docker/baseline/assets/icons
	cp js/dist/icons/grocery.svg docker/baseline/assets/icons
	cp js/dist/icons/school.svg docker/baseline/assets/icons
	cp js/dist/icons/childcare.svg docker/baseline/assets/icons
	cp js/dist/icons/restaurant.svg docker/baseline/assets/icons
	cp js/dist/icons/walking.svg docker/baseline/assets/icons
	cp js/dist/icons/cycle.svg docker/baseline/assets/icons
	cp js/dist/icons/rail.svg docker/baseline/assets/icons
	cp js/dist/icons/town.svg docker/baseline/assets/icons
	cp js/dist/icons/vacant.svg docker/baseline/assets/icons
	cp js/dist/icons/select-arrow.svg docker/baseline/assets/icons
	cp js/dist/icons/select-missing.svg docker/baseline/assets/icons
	cp js/dist/icons/run-arrow.svg docker/baseline/assets/icons
	cp js/dist/icons/run-arrow-hover.svg docker/baseline/assets/icons
	mkdir -p docker/baseline/static
	cp src/diagonal.works/diagonal/cmd/baseline/bundle.js docker/baseline/static
	cp src/diagonal.works/diagonal/cmd/baseline/main.css docker/baseline/static
	cp src/diagonal.works/diagonal/cmd/baseline/index.html docker/baseline/static
	mkdir -p docker/baseline/data
	cp src/diagonal.works/diagonal/cmd/baseline/galashiels.geojson docker/baseline/data
	docker build --build-arg platform=${TARGETPLATFORM} -f docker/Dockerfile.baseline -t baseline-${TARGETARCH} docker
	docker tag baseline-${TARGETARCH} eu.gcr.io/diagonal-platform/baseline-${TARGETARCH}
	docker push eu.gcr.io/diagonal-platform/baseline-${TARGETARCH}

# Use TARGETARCH=x86_64 TARGETOS=linux for GCP
docker-b6:
	rm -rf docker/b6
	mkdir -p docker/bin/${TARGETPLATFORM}
	cp bin/${TARGETPLATFORM}/b6 docker/bin/${TARGETPLATFORM}
	mkdir -p docker/b6/js
	cp src/diagonal.works/b6/cmd/b6/js/bundle.js docker/b6/js
	cp -r src/diagonal.works/b6/cmd/b6/js/static docker/b6/
	docker build --build-arg platform=${TARGETPLATFORM} -f docker/Dockerfile.b6 -t b6-${TARGETARCH} docker
	docker tag b6-${TARGETARCH} eu.gcr.io/diagonal-platform/b6-${TARGETARCH}
	docker push eu.gcr.io/diagonal-platform/b6-${TARGETARCH}

protos:
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=proto --go_out=src proto/cookie.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=proto --go_out=src proto/tiles.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=proto --go_out=src proto/osm.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=proto --go_out=src proto/geometry.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=proto --go_out=src proto/features.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go --plugin=${HOME}/go/bin/protoc-gen-go-grpc -I=proto --go_out=src --go-grpc_out=src proto/api.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=src/diagonal.works/diagonal/osm --go_out=src src/diagonal.works/diagonal/osm/import.proto
	protoc --plugin=${HOME}/go/bin/protoc-gen-go -I=src/diagonal.works/b6/osm/proto --go_out=src src/diagonal.works/b6/osm/proto/pbf.proto

src/diagonal.works/b6/api/y.go: src/diagonal.works/b6/api/shell.y
	cd src/diagonal.works/b6/api; goyacc shell.y

experimental:
	cd src/diagonal.works/diagonal/experimental/mr; go build
	cd src/diagonal.works/diagonal/experimental/osmpbf; go build

experimental_gazetteer:
	cd src/diagonal.works/diagonal/experimental/gazetteer; go build

experimental_sightline_tiles:
	cd src/diagonal.works/diagonal/monitoring; go generate
	cd src/diagonal.works/diagonal/experimental/sightline-tiles; go build -o ../../../../../bin/${TARGETPLATFORM}/sightline-tiles

experimental_pyramid_tiles:
	cd src/diagonal.works/diagonal/monitoring; go generate
	cd src/diagonal.works/diagonal/experimental/pyramid-tiles; go build -o ../../../../../bin/${TARGETPLATFORM}/pyramid-tiles

experimental_posting_lists:
	cd src/diagonal.works/diagonal/experimental/posting-lists; go build -o ../../../../../bin/${TARGETPLATFORM}/posting-lists

experimental_s2-sharding:
	cd src/diagonal.works/diagonal/experimental/s2-sharding; go build -o ../../../../../bin/${TARGETPLATFORM}/s2-sharding

experimental_collada:
	cd src/diagonal.works/diagonal/experimental/collada; go build -o ../../../../../bin/${TARGETPLATFORM}/collada

experimental_coding_exercise:
	cd src/diagonal.works/diagonal/experimental/coding-exercise; go build -o ../../../../../bin/${TARGETPLATFORM}/coding-exercise

python: python/diagonal_b6/api_generated.py

python/diagonal_b6/api_generated.py: b6-api
	bin/${TARGETPLATFORM}/b6-api | python/diagonal_b6/generate_api.py > $@
	python3 -m grpc_tools.protoc -Iproto --python_out=python/diagonal_b6 --grpc_python_out=python/diagonal_b6  proto/geometry.proto proto/features.proto proto/api.proto
	sed -e 's/import geometry_pb2/import diagonal_b6.geometry_pb2/' -i "" python/diagonal_b6/features_pb2.py
	sed -e 's/import geometry_pb2/import diagonal_b6.geometry_pb2/' -i "" python/diagonal_b6/api_pb2.py
	sed -e 's/import features_pb2/import diagonal_b6.features_pb2/' -i "" python/diagonal_b6/api_pb2.py
	sed -e 's/import api_pb2/import diagonal_b6.api_pb2/' -i "" python/diagonal_b6/api_pb2_grpc.py

ipython: python
	cd python; pip3 install . --upgrade --target ${HOME}/.ipython/

python-test: python b6
	PYTHONPATH=python TARGETPLATFORM=${TARGETPLATFORM} python3 python/diagonal_b6/b6_test.py

test:
	cd src/diagonal.works/b6; go test diagonal.works/b6/...
	cd src/diagonal.works/diagonal; go test diagonal.works/diagonal/...

clean:
	find . -type f -perm +a+x | xargs rm

.PHONY: python
