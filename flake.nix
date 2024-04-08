{
  description = "DigitalService GmbH des Bundes - Platform CLI";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: let
    bin = pkgs: pkgs.buildGo122Module rec {
      name = "git-commit-template";
      src = ./.;
      CGO_ENABLED = 0;
      vendorHash = pkgs.lib.fileContents ./go.mod.sri;
      nativeCheckInputs = [pkgs.git];
    };
    flakeForSystem = nixpkgs: system: let
      pkgs = nixpkgs.legacyPackages.${system};
      git-commit-template = bin pkgs;
    in {
      packages = {
        default = git-commit-template;
      };
    };
  in
    flake-utils.lib.eachDefaultSystem (system: flakeForSystem nixpkgs system);
}
