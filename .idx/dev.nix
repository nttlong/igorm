# To learn more about how to use Nix to configure your environment
# see: https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "stable-24.05"; # or "unstable"
  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
    pkgs.air
  ];
  # Sets environment variables in the workspace
  env = {};
  idx = {
    previews = {
      enable = true;
      previews = {
        web = {
          #dir = "/home/user/test/igorm/examples/media"; # set working directory
          # Example: run "npm run dev" with PORT set to IDX's defined port for previews,
          # and show it in IDX's web preview panel
          
          
          command = [ "sh" "-c" "go run cmd/main.go" ];
          manager = "web";
          env = {
            # Environment variables to set for your server
            PORT = "$PORT";
          };
        };
      };
    };
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
    ];
    workspace = {
      onCreate = {
        # Open editors for the following files by default, if they exist:
        #default.openFiles = ["main.go"];
      };
      # Runs when a workspace is (re)started
      onStart= {
        #run-server = "air";
      };
      # To run something each time the workspace is first created, use the `onStart` hook
    };
  };
}
