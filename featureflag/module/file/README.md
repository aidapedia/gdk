# File-based Feature Flag Module

This module provides a file-based implementation for the feature flag system. It allows you to load feature flags from a JSON file.

## JSON Structure

The feature flags are defined in a JSON file. The structure supports nested objects, which are treated as folders.

Example `config.json`:

```json
{
    "key_bool": true,
    "key_string": "some string value",
    "key_int": 42,
    "key_float": 3.14,
    "feature_group": {
        "sub_feature_enabled": true,
        "config_value": "nested value"
    }
}
```

## Usage

### Initialization

To use the file module, you need to initialize it with the path to your JSON configuration file. You can optionally provide a prefix to scope the feature flags to a specific section of the JSON.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aidapedia/gdk/featureflag/module/file"
)

func main() {
	// Initialize with the path to the JSON file and an empty prefix
	ffModule := file.New("/path/to/config.json", "")

	// Or initialize with a prefix to scope to "feature_group"
	// ffModule := file.New("/path/to/config.json", "feature_group")
}
```

### Retrieving Values

You can retrieve values using the standard methods provided by the module interface. Note that for nested keys, you should use the `/` separator.

```go
ctx := context.Background()

// Get a boolean value
enabled, err := ffModule.GetBool(ctx, "key_bool")
if err != nil {
    log.Printf("Error getting bool: %v", err)
}

// Get a string value
strVal, err := ffModule.GetString(ctx, "key_string")
if err != nil {
    log.Printf("Error getting string: %v", err)
}

// Get an integer value
intVal, err := ffModule.GetInt(ctx, "key_int")
if err != nil {
    log.Printf("Error getting int: %v", err)
}

// Get a value from a nested structure (if initialized without prefix)
// Use slash separator for nested keys
nestedVal, err := ffModule.GetBool(ctx, "feature_group/sub_feature_enabled") 

// Get a struct value
// The value in JSON must be a string containing the JSON representation of the struct
// Example JSON: "key_struct": "{\"field\": \"value\"}"
type MyStruct struct {
    Field string `json:"field"`
}
var myStruct MyStruct
err = ffModule.GetStruct(ctx, "key_json_string", &myStruct)
if err != nil {
    log.Printf("Error getting struct: %v", err)
}
```

### Watching for Changes

The `Watch` method allows you to monitor the configuration file for changes. It polls the file every 5 seconds. If a change is detected, it updates the internal state and sends a signal on the returned channel.

```go
// Start watching for changes
changeChan, err := ffModule.Watch(ctx)
if err != nil {
    log.Fatalf("Failed to watch for changes: %v", err)
}

// Listen for changes in a goroutine
go func() {
    for range changeChan {
        log.Println("Configuration file changed! Reloading values...")
        // You can now retrieve the updated values
        newVal, _ := ffModule.GetBool(ctx, "key_bool")
        log.Printf("New value for key_bool: %v", newVal)
    }
}()
```
