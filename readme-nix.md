# Notes on using this repo with Nix

Overall, the project builds well with Nix; the main issue is that a couple of
the Makefile targets don't work at present:

- [ ] The `version` calculation seems to not work on a branch (?? to confirm)
- [ ] An `npm install -g pnpm` fails on NixOS

Note that all commands below assume they are executed from the root directory
of the repo.

## Building the code

### Building/using the go binaries

```shell
nix build
```

Binaries are then located at:

```shell
./result/bin/b6
./result/bin/b6-ingest-osm
... # and so on for all the binaries
```

### Building the UI (v2)

You need to be in the devShell:

```shell
nix develop
```

Then, to setup:

```shell
cd frontend
pnpm install --config.confirmModulesPurge=false
```

Then, to build:

```shell
pnpm build
```

or to run a dev server:

```shell
pnpm dev
```

### Building the UI (v1)

You don't need to do this to run the V2 UI; but if you want to have it:

```shell
cd /src/diagonal.works/b6/cmd/b6/js
make
```

## Compute an index (ingesting data)

Follows the readme, just uses the Nix binary:

```shell
./result/bin/b6-ingest-osm --input=data/tests/camden.osm.pbf --output=data/camden.index
```

## Running

Note that the above has only built the V2 UI; so we can run b6 with that UI as
follows:

```shell
./result/bin/b6 -world data/camden.index -enable-v2-ui
```

## Obtaining data

One method to obtain data is to download the entire planet data from OSM:

```shell
docker run --rm -it -v \
         $PWD:/download openmaptiles/openmaptiles-tools \
         download-osm \
         planet -- -d /download
```

Which for me resulted in a 76 gb file.

From there, you can use `osmium` (included in the nix development shell) to
extract specific regions.

I grabbed polygons from here: <https://download.openstreetmap.fr/polygons/>
and then ran the following:

```shell
osmium extract \
         --polygon ~/tmp/osm/victoria.poly \
         -o victoria.osm.pbf \
         -F pbf \
         ~/tmp/osm/planet-240708.osm.pbf
```

to extract the Victoria polygon from the data, using the planet download I had
in `~/tmp/osm` and the `.poly` file.

This then gives a `.osm.pbf` file that can be ingested using the earlier
method.

For more, see also the
[download-osm](https://github.com/openmaptiles/openmaptiles-tools#multi-streamed-osm-data-downloader)
tool docs directly.


### Trivia

- If you use the `=` style, the CLI doesn't allow for shell-based path
  expansion; so i.e. `b6 -world=~/dev/diagonal/b6/data/camden.index` would fail,
  but `b6 -world ~/dev/diagonal/b6/data/camden.index` works fine.
- You will need to re-run gomod2nix, I think, if new go dependencies are
  added? (to confirm)
- Because `go` is also in the devShell, you can build the project there as
  well (note: actually, you can't, because of this `VERSION` problem above, at
  the moment.)

### Todo

- [ ] The tests fail, so the flake has `doCheck = false`. I haven't worked out
  how to run the tests outside of that, yet.



### Random things to investigate

- [ ] We could use a `not` query (like `or`) for filtering
- [ ] It's a bit weird to use `.all()`; does that have a performance impact?
- [ ] Build the Makefile items as a nix derivation?
- [ ] Document go <-> Python api (with types!) somewhere
- [x] Get the tests running
- [ ] Add some tests to highlight bad functionality
      - Not being able to filter invalid
        ```
        b6.filter_accessible(q, b6.all())
        ```
        resulting in less things.
- [x] Add a `drop_invalid` command to drop based on feature-id invalid
