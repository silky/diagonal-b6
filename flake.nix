{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/24.05";
    unstable.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    pyproject-nix.url = "github:nix-community/pyproject.nix";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };


  outputs =
    { self
    , nixpkgs
    , gomod2nix
    , unstable
    , pyproject-nix
    , flake-utils
    , treefmt-nix
    , ...
    }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      # Python setup
      overlay = _: prev: {
        python3 = prev.python3.override {
          packageOverrides = _: p: {
            s2sphere = p.buildPythonPackage rec {
              version = "0.2.5";
              pname = "s2sphere";
              format = "pyproject";
              nativeBuildInputs = with p.pythonPackages; [
                setuptools
              ];
              propagatedBuildInputs = with p.pythonPackages; [
                future
              ];
              src = pkgs.fetchFromGitHub {
                owner = "silky";
                repo = "s2sphere";
                rev = "d1d067e8c06e5fbaf0cc0158bade947b4a03a438";
                sha256 = "sha256-6hNIuyLTcGcXpLflw2ajCOjel0IaZSFRlPFi81Z5LUo=";
              };
            };
          };
        };
      };

      python = pkgs.python3;

      pythonProject = pyproject-nix.lib.project.loadPyproject {
        projectRoot = ./python;
      };

      pythonEnv = python.withPackages (ps:
          pythonProject.renderers.withPackages {
            inherit python;
          } ps ++
          [
            # For `make python`
            ps.grpcio-tools
            ps.jupyter
          ]
      );

      b6-py = python.pkgs.buildPythonPackage
        (pythonProject.renderers.buildPythonPackage {
          inherit python;
        });

      # Go setup
      b6-go = with pkgs; gomod2nix.legacyPackages.${system}.buildGoApplication {
        name = "b6";
        src = ./src/diagonal.works/b6;
        buildInputs = [
          gdal
        ];
        nativeBuildInputs = [
          pkg-config
        ];

        # The tests fail presently
        doCheck = false;

        # Must be added due to bug https://github.com/nix-community/gomod2nix/issues/120
        pwd = ./src/diagonal.works/b6;

        # Optional flags.
        # CGO_ENABLED = 0;
        # flags = [ "-trimpath" ];
        # ldflags = [ "-s" "-w" "-extldflags -static" ];
      };

      pkgs = import nixpkgs { inherit system; overlays = [ overlay ]; };
      unstablePkgs = import unstable { inherit system; };
    in
    {
      # Development shells for hacking
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          # Running the Makefile tasks
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
          pkg-config gdal

          # Python hacking
          pythonEnv

          # Go hacking
          go_1_21
          gotools
          gomod2nix.packages.${system}.default # gomod2nix CLI

          # JS Hacking
          nodejs
          unstablePkgs.pnpm # Need version 9

          # Other
          osmium-tool # Extract OSM
        ];

        shellHook = ''
          export PYTHONPATH=''$(pwd)/python
        '';
      };


      packages = {
        default = b6-go;
        go = b6-go;
        python = b6-py;
      };


      formatter =
        let
          fmt = treefmt-nix.lib.evalModule pkgs (_: {
            projectRootFile = "flake.nix";
            programs.nixpkgs-fmt.enable = true;
          });
        in
        fmt.config.build.wrapper;
    });
}
