# WASM Plugin System - Proof of Concept

## Overview

PicoClaw WASM Plugin System cho phép mở rộng chức năng thông qua WebAssembly plugins sử dụng `wazero` runtime.

## Architecture

```
┌─────────────────────────────────────┐
│     PicoClaw Core Application       │
├─────────────────────────────────────┤
│       Plugin Manager (Go)           │
│  - Load/Unload plugins              │
│  - Manage plugin lifecycle          │
│  - Handle plugin metadata           │
├─────────────────────────────────────┤
│     WASM Runtime (wazero)           │
│  - Execute WASM modules             │
│  - Provide WASI support             │
│  - Memory isolation                 │
├─────────────────────────────────────┤
│         WASM Plugins                │
│  - Written in any WASM language     │
│  - Sandboxed execution              │
│  - Export functions to host         │
└─────────────────────────────────────┘
```

## Features

- ✅ **Sandboxed Execution**: Plugins run in isolated WASM environment
- ✅ **Multi-language Support**: Write plugins in Go (TinyGo), Rust, C, etc.
- ✅ **Hot Reload**: Load/unload plugins at runtime
- ✅ **Metadata Management**: Track plugin versions and capabilities
- ✅ **WASI Support**: Basic I/O operations available

## Quick Start

### 1. Create a Plugin (TinyGo)

```go
//go:build tinygo

package main

//export add
func add(a, b int32) int32 {
    return a + b
}

//export process
func process(input int32) int32 {
    return input * 2
}

func main() {}
```

### 2. Compile to WASM

```bash
tinygo build -o plugin.wasm -target=wasi plugin.go
```

### 3. Load and Use Plugin

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/sipeed/picoclaw/pkg/plugins"
)

func main() {
    ctx := context.Background()
    
    // Create plugin manager
    pm, err := plugins.NewPluginManager(ctx)
    if err != nil {
        panic(err)
    }
    defer pm.Close(ctx)
    
    // Read WASM file
    wasmBytes, err := os.ReadFile("plugin.wasm")
    if err != nil {
        panic(err)
    }
    
    // Load plugin
    metadata := &plugins.PluginMetadata{
        Name:        "math_plugin",
        Version:     "1.0.0",
        Description: "Simple math operations",
        Author:      "PicoClaw Team",
        Capabilities: []string{"add", "process"},
    }
    
    err = pm.LoadPlugin(ctx, metadata, wasmBytes)
    if err != nil {
        panic(err)
    }
    
    // Call plugin function
    result, err := pm.runtime.CallFunction(ctx, "math_plugin", "add", 5, 7)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("5 + 7 = %d\n", result[0])
}
```

## API Reference

### WASMRuntime

```go
// Create new runtime
runtime, err := plugins.NewWASMRuntime(ctx)

// Load plugin
err = runtime.LoadPlugin(ctx, "plugin_name", wasmBytes)

// Call function
result, err := runtime.CallFunction(ctx, "plugin_name", "function_name", args...)

// Unload plugin
err = runtime.UnloadPlugin(ctx, "plugin_name")

// Close runtime
err = runtime.Close(ctx)
```

### PluginManager

```go
// Create manager
pm, err := plugins.NewPluginManager(ctx)

// Load with metadata
err = pm.LoadPlugin(ctx, metadata, wasmBytes)

// Get plugin info
meta, exists := pm.GetPlugin("plugin_name")

// List all plugins
list := pm.ListPlugins()

// Unload plugin
err = pm.UnloadPlugin(ctx, "plugin_name")
```

## Testing

Run tests:

```bash
go test ./pkg/plugins/... -v
```

Expected output:
```
=== RUN   TestWASMRuntimeCreation
--- PASS: TestWASMRuntimeCreation (0.00s)
=== RUN   TestPluginManagerCreation
--- PASS: TestPluginManagerCreation (0.00s)
=== RUN   TestSimpleWASMPlugin
--- PASS: TestSimpleWASMPlugin (0.01s)
=== RUN   TestPluginManagerWithMetadata
--- PASS: TestPluginManagerWithMetadata (0.01s)
PASS
```

## Example Plugins

### 1. Math Plugin (TinyGo)

```go
//go:build tinygo

package main

//export add
func add(a, b int32) int32 {
    return a + b
}

//export subtract
func subtract(a, b int32) int32 {
    return a - b
}

//export multiply
func multiply(a, b int32) int32 {
    return a * b
}

func main() {}
```

### 2. String Plugin (Rust)

```rust
#[no_mangle]
pub extern "C" fn string_length(ptr: *const u8, len: usize) -> usize {
    let slice = unsafe { std::slice::from_raw_parts(ptr, len) };
    slice.len()
}

#[no_mangle]
pub extern "C" fn to_uppercase(ptr: *mut u8, len: usize) {
    let slice = unsafe { std::slice::from_raw_parts_mut(ptr, len) };
    for byte in slice {
        if *byte >= b'a' && *byte <= b'z' {
            *byte -= 32;
        }
    }
}
```

## Security Considerations

1. **Sandboxing**: WASM provides memory isolation
2. **Resource Limits**: Set execution timeouts
3. **Capability Control**: Limit WASI capabilities
4. **Input Validation**: Validate all plugin inputs
5. **Version Control**: Track plugin versions

## Performance

- **Startup**: ~1-5ms per plugin load
- **Execution**: Near-native performance
- **Memory**: Isolated per plugin
- **Overhead**: Minimal (~100KB per runtime)

## Roadmap

- [ ] Memory sharing between host and plugins
- [ ] Plugin configuration system
- [ ] Plugin dependency management
- [ ] Hot reload without restart
- [ ] Plugin marketplace/registry
- [ ] Advanced WASI features
- [ ] Plugin communication protocol

## Limitations

- No direct access to Go runtime
- Limited standard library in TinyGo
- Memory must be managed carefully
- No threading support yet

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## License

MIT License - See [LICENSE](../../LICENSE)
