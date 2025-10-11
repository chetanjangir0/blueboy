{
  description = "Blueboy - A simple TUI network manager";
  
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "blueboy";
          version = "0.1.0";
          
          src = ./.;
          
          # Hash of Go module dependencies
          vendorHash = "sha256-3RgtCmWae1OqD3hBO84WghW7rlZ6yRLw9WHRIvq/qsY=";
          
          # Use proxy vendor to avoid vendor directory issues
          proxyVendor = true;
          
          # Build from cmd/tui directory
          subPackages = [ "cmd/tui" ];
          
          # Runtime dependency on nmcli (NetworkManager)
          buildInputs = [ pkgs.networkmanager ];
          
          nativeBuildInputs = [ pkgs.makeWrapper ];
          
          # Ensure nmcli is available at runtime
          postInstall = ''
            wrapProgram $out/bin/blueboy \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.networkmanager ]}
          '';
          
          # Optional: embed version info
          ldflags = [
            "-s"
            "-w"
            "-X main.version=${self.rev or "dev"}"
          ];
          
          meta = with pkgs.lib; {
            description = "A simple TUI network manager";
            homepage = "https://github.com/chetanjangir0/blueboy";
            license = licenses.mit;
            maintainers = [ "chetan" ];
            mainProgram = "tui";
            platforms = platforms.linux;
          };
        };
        
        # Development shell with all necessary tools
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
        
        # Optional: create an app output for easy running
        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/tui";
        };
      }
    );
}
