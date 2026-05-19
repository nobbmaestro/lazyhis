{
  lib,
  buildGoModule,
  pkg-config,
  stdenv,
  libx11,
  libxi,
  libxfixes,
  src,
  gitRevision ? "dev",
}:
buildGoModule rec {
  pname = "lazyhis";
  version = "0.10.0";

  inherit src;

  vendorHash = "sha256-8tQB9rQfk5iy5dJ6n8RemKiFIf5bToYXhLWWxx/y+dM=";

  doCheck = true;

  nativeBuildInputs = [
    pkg-config
  ];

  buildInputs = lib.optionals stdenv.isLinux [
    libx11
    libxi
    libxfixes
  ];

  ldflags = [
    "-X main.version=${version}"
    "-X main.commit=${gitRevision}"
  ];

  postInstall = ''
    export HOME=$(mktemp -d)
    mkdir -p $out/share/man/man1
    $out/bin/lazyhis gen man --dst $out/share/man/man1
  '';

  meta = with lib; {
    description = "A simple terminal UI for shell history, written in Go!";
    homepage = "https://github.com/nobbmaestro/lazyhit";
    license = licenses.mit;
    maintainers = [ "nobbmaestro" ];
    platforms = platforms.unix;
    mainProgram = "lazyhis";
  };
}
