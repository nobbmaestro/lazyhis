{
  description = "A simple terminal UI for shell history, written in Go!";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    systems.url = "github:nix-systems/default";
  };

  outputs =
    inputs@{
      flake-parts,
      systems,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      perSystem =
        {
          pkgs,
          system,
          ...
        }:
        let
          lazyhis = pkgs.buildGoModule rec {
            pname = "lazyhis";
            version = "0.9.6";

            gitCommit = inputs.self.rev or inputs.self.dirtyRev or "dev";

            src = inputs.self;

            vendorHash = "sha256-8tQB9rQfk5iy5dJ6n8RemKiFIf5bToYXhLWWxx/y+dM=";

            doCheck = true;

            ldflags = [
              "-X main.version=${version}"
              "-X main.commit=${gitCommit}"
            ];

            postInstall = ''
              export HOME=$(mktemp -d) 
              mkdir -p $out/share/man/man1 
              $out/bin/lazyhis gen man --dst $out/share/man/man1 
            '';

            meta = {
              description = "A simple terminal UI for shell history, written in Go!";
              homepage = "https://github.com/nobbmaestro/lazyhit";
              license = pkgs.lib.licenses.mit;
              maintainers = [ "nobbmaestro" ];
              platforms = pkgs.lib.platforms.unix;
              mainProgram = "lazyhis";
            };
          };
        in
        {
          packages = {
            default = lazyhis;
            inherit lazyhis;
          };

          devShells.default = pkgs.mkShell {
            name = "lazyhis-dev";

            packages = [
              pkgs.go
              pkgs.golangci-lint
              pkgs.golines
              pkgs.gnumake
            ];
          };
      flake = {
        overlays.default = final: prev: {
          lazyhis = inputs.self.packages.${final.system}.lazyhis;
        };
      };
    };
}
