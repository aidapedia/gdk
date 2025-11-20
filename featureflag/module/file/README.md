# Feature Flag

The `featureflag` package provides a flexible way to manage feature flags in your application. It supports different modules, including a file-based module.

## JSON Structure Format

When using the **File Module**, the configuration is stored in a JSON file. The structure consists of a `root` object which contains `keys` for values and `child` for nested configurations.

> [!IMPORTANT]
> The JSON tag for nested directories is `child` (singular), not `children`.

### Example `config.json`

```json
{
  "root": {
    "keys": {
      "global_feature_enabled": true,
      "max_retries": 3,
      "app_name": "MyApp",
      "config_json": "{\"timeout\": 5000}"
    },
    "child": {
      "payment": {
        "keys": {
          "stripe_enabled": true
        },
        "child": {
          "v2": {
            "keys": {
              "new_flow": true
            }
          }
        }
      }
    }
  }
}
```

## Initialization

To initialize the feature flag with the file module:

```go
package main

import (
	"github.com/aidapedia/gdk/featureflag"
	"github.com/aidapedia/gdk/featureflag/module"
)

func main() {
	// Initialize the feature flag
	ff := featureflag.New(featureflag.Option{
		Module:  module.FileModule,
		Address: "path/to/config.json", // Path to your JSON file
		Prefix:  "/",                   // Optional prefix to scope the root
	})

    // ... use ff
}
```

## Getting Keys

You can retrieve values using various typed methods. All methods accept a `context.Context` and the key path.

```go
ctx := context.Background()

// Get a boolean value
enabled, err := ff.GetBool(ctx, "global_feature_enabled")

// Get an integer value
retries, err := ff.GetInt(ctx, "max_retries")

// Get a string value
appName, err := ff.GetString(ctx, "app_name")

// Get a value from a nested path (separated by /)
stripeEnabled, err := ff.GetBool(ctx, "payment/stripe_enabled")

// Get a raw interface{} value
val, err := ff.GetValue(ctx, "some/key")

// Get a struct (parses JSON string value)
// Note: The value in the JSON file must be a string containing JSON.
type MyConfig struct {
    Timeout int `json:"timeout"`
}
var cfg MyConfig
err = ff.GetStruct(ctx, "config_json", &cfg)
```

## Watching Feature Flag Values

You can watch for changes in the feature flag configuration file. The `Watch` method returns a channel that receives a notification whenever the file changes.

```go
// Start watching
watcher, err := ff.Watch(ctx)
if err != nil {
    // handle error
}

go func() {
    for range watcher {
        // The feature flag configuration has changed
        // You can re-fetch values or trigger updates here
        fmt.Println("Feature flag configuration updated")
        
        // Example: re-check a flag
        newVal, _ := ff.GetBool(ctx, "global_feature_enabled")
        fmt.Println("New value:", newVal)
    }
}()
```
