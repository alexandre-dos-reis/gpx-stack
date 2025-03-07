{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  pkgs-unstable = import inputs.nixpkgs-unstable {system = pkgs.stdenv.system;};

  DB_OWNER = "go";
  DB_PASSWORD = "lang";
  DB_NAME = DB_OWNER;
  DB_PORT = 8786;
  DB_HOST = "localhost";
  DB_URL = "postgresql://${DB_OWNER}:${DB_PASSWORD}@${DB_HOST}:${toString DB_PORT}/${DB_NAME}?sslmode=disable";
in {
  env.DB_URL = DB_URL;

  packages = with pkgs; [
    air
    sqlc
    go-migrate
    bun
  ];

  services.postgres = {
    enable = true;
    port = DB_PORT;
    listen_addresses = DB_HOST;
    package = pkgs.postgresql_17;
    initialDatabases = [
      {
        name = DB_OWNER;
        user = DB_OWNER;
        pass = DB_PASSWORD;
      }
    ];
  };

  languages.go = {
    enable = true;
    package = pkgs-unstable.go;
  };

  processes = {
    air.exec = "air";
    vite.exec = "bun dev";
  };

  scripts = let
    migrationPath = "database/migrations";
  in {
    mg-create.exec = ''
      migrate create -ext sql -dir ${migrationPath} -seq $1
    '';
    mg-migrate-up.exec = ''
      migrate -database ${DB_URL} -path ${migrationPath} up
      sqlc generate
    '';
    mg-migrate-down.exec = ''
      migrate -database ${DB_URL} -path ${migrationPath} down $1
      sqlc generate
    '';
    sql.exec = ''
      psql ${DB_URL}

    '';
  };
}
