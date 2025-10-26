{
  description = "Blueboy - A simple TUI network manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "blueboy";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-3RgtCmWae1OqD3hBO84WghW7rlZ6yRLw9WHRIvq/qsY=";
          proxyVendor = true;
          buildInputs = [ pkgs.networkmanager ];
          nativeBuildInputs = [ pkgs.makeWrapper ];

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${self.rev or "dev"}"
          ];

          meta = with pkgs.lib; {
            description = "A simple TUI network manager";
            homepage = "https://github.com/chetanjangir0/blueboy";
            license = licenses.mit;
            maintainers = [ "chetanjangir0" ];
            mainProgram = "blueboy";
            platforms = platforms.linux;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            networkmanager
          ];

          shellHook = ''
            echo "Blueboy development environment"
            echo "Go version: $(go version)"
            echo "nmcli available: $(which nmcli)"
          '';
        };

      }
    );
}
