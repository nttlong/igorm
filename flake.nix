{
  description = "Dev environment with Go and Git";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";

  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      devShells.default = pkgs.mkShell {
        buildInputs = [
          pkgs.go
          pkgs.git
        ];
      };
    };
}
