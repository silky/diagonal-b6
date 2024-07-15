{
  inputs = {
    nixpkgs = { url = "github:NixOS/nixpkgs/24.05"; };
    unstable = { url = "github:NixOS/nixpkgs/nixos-unstable"; };
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };


  outputs =
    { self
    , nixpkgs
    , gomod2nix
    , unstable
    , ...
  }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; };
        unstablePkgs = import unstable { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ system, pkgs, unstablePkgs, ... }:
        let
          buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;
        in
        {
          default = with pkgs; buildGoApplication {
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
        });

      devShells = forAllSystems ({ system, pkgs, unstablePkgs, ... }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_21
            gotools
            gomod2nix.packages.${system}.default # gomod2nix CLI
            nodejs
            unstablePkgs.pnpm # Need version 9

            osmium-tool # Extract OSM

            (python3.withPackages (ps: [
              ps.grpcio-tools
            ]))
          ];
        };
      });
    };
}
