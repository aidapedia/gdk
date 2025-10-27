# Config Manager

This package provides a simple way to load configuration and secrets into your application using `viper`. It supports merging multiple config files and reading secrets from either a file or Google Secret Manager (GSM).

## Overview

- Load config from one or more files in a directory.
- Choose a top-level key to unmarshal from (e.g. `Config`).
- Read secrets from a file path or GSM, controlled via environment variables.

## Environment Variables

- `CONFIG_FILE_PATH`: Directory containing your config files (e.g. `./config` or `./config/test`).
- `SECRET_FILE_PATH`: Full path to the secret file (e.g. `./config/secret.json`).
- `SECRET_GSM_PROJECT_ID`: GSM project ID when using GSM.

## Config File Layout

You can split config across multiple files. Each file should have a top-level key that matches the `key` you pass to the manager (commonly `Config`).

Example directory structure:

```
config/
  config.json
  app.json
```

Example file contents:

`config/config.json`
```json
{
  "Config": {
    "Environment": "dev"
  }
}
```

`config/app.json`
```json
{
  "Config": {
    "AppEnv": "test"
  }
}
```

## Quick Start (Load Config)

```go
package main

import (
    "context"
    gdkconfig "github.com/aidapedia/gdk/config"
)

// Define your target struct that matches the JSON inside the top-level key.
type AppConfig struct {
    Environment string
    AppEnv      string
}

func main() {
    // Ensure environment variable points to the directory with your config files.
    // e.g. os.Setenv("CONFIG_FILE_PATH", "./config")

    cfg := AppConfig{}

    // file names are specified WITHOUT extensions and must exist under CONFIG_FILE_PATH.
    // key is the top-level JSON key inside each file (e.g. "Config").
    m := gdkconfig.NewManager(Option{
				TargetStore: &cfg,
				FileName:    []string{"config", "app"},
				ConfigKey:   "AppConfig",
			})

    if err := m.SetConfig(context.Background()); err != nil {
        panic(err)
    }

    // Get your config imediately
    _ = cfg
}
```

### Notes
- Pass your target struct pointer (`cfg`) as the first argument to `NewManager`.
- The manager will `MergeInConfig()` for each file name provided and `UnmarshalKey(key, &store)`.
- If you want to unmarshal the whole file, put your fields under a common top-level key and use that key.

## Load Secrets

You can also load secrets using the manager or the `secret` package directly.

### Secrets via Manager (File)
```go
package main

import (
    "context"
    gdkconfig "github.com/aidapedia/gdk/config"
)

type Secrets struct {
    APIKey string
    Token  string
}

func main() {
    // SECRET_FILE_PATH must be set to the full file path, e.g. ./config/secret.json
    // e.g. os.Setenv("SECRET_FILE_PATH", "./config/secret.json")

    sec := &Secrets{}
    m := gdkconfig.NewManager(nil, gdkconfig.SecretTypeFile, nil, "")

    // Internally, SetSecretStore reads SECRET_FILE_PATH and fills m.secretStore.
    // To use it, set m.secretStore to your target type first.
    // (Do this inside the same package or provide a wrapper that sets it.)
}
```

### Secrets via `secret.File` (Direct)
If you prefer to load secrets directly:
```go
package main

import (
    "context"
    "github.com/aidapedia/gdk/config/secret"
)

type Secrets struct {
    ServiceName string
}

func main() {
    // Provide the full path to the secret file.
    // If your file is in a directory, include the directory in the path.
    s := secret.NewSecretFile("config/secret.json")

    var sec Secrets
    if err := s.GetSecret(context.Background(), &sec); err != nil {
        panic(err)
    }
}
```

Example `config/secret.json`:
```json
{
  "ServiceName": "example-service"
}
```

## Troubleshooting

- "CONFIG_FILE_PATH environment variable is not set": Set `CONFIG_FILE_PATH` to the directory that contains your config files.
- "Failed to get secret: open <file>: no such file or directory": Ensure you pass the full relative/absolute path to the secret file (`secret.NewSecretFile("config/secret.json")`) or set `SECRET_FILE_PATH` correctly.
- Using relative paths: Paths are resolved from the process working directory. In tests inside `config`, setting `CONFIG_FILE_PATH` to `"test"` reads from `config/test`.

## Testing Example

The repository includes tests under `config/test/` demonstrating multi-file merge with the `Config` key:
- `config/test/config.json` sets `Environment`.
- `config/test/app.json` sets `AppEnv`.
- The test sets `CONFIG_FILE_PATH` to `"test"` and uses file names `[]string{"config", "app"}`.

This pattern lets you compose configuration cleanly across files and environments.