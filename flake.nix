{
  description = "A simple terminal UI for shell history, written in Go!";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    systems.url = "github:nix-systems/default";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    inputs@{
      flake-parts,
      systems,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        { system, ... }:
        let
          goMod = builtins.readFile ./go.mod;
          versionMatch = builtins.match ".*go[[:space:]]([0-9]+\\.[0-9]+)(\\.[0-9]+)?.*" goMod;

          goVersion =
            if versionMatch != null then
              builtins.head versionMatch
            else
              throw "Could not extract Go version from go.mod";

          goOverlay = final: prev: {
            go = prev."go_${builtins.replaceStrings [ "." ] [ "_" ] goVersion}";
          };

          pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [ goOverlay ];
            config = { };
          };

          lazyhis = pkgs.callPackage ./default.nix {
            src = inputs.self;
            gitRevision = inputs.self.rev or inputs.self.dirtyRev or "dev";
          };
        in
        {
          _module.args.pkgs = pkgs;

          packages = {
            default = lazyhis;
            inherit lazyhis;
          };

          apps.default = {
            type = "app";
            program = "${lazyhis}/bin/lazyhis";
          };

          devShells.default = pkgs.mkShell {
            name = "lazyhis-dev";

            packages = [
              pkgs.go
              pkgs.gotools
              pkgs.golangci-lint
              pkgs.golines
              pkgs.gnumake
            ];
          };

          treefmt = {
            programs.nixfmt.enable = true;
            programs.gofmt.enable = true;
          };

          checks.build = lazyhis;
        };

      flake = {
        overlays.default = final: prev: {
          lazyhis = inputs.self.packages.${final.system}.lazyhis;
        };
      };
    };
}
