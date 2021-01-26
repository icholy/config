# Config Format

> This package implements an unmarshaler for a simplified version of the HCL configuration language. The main difference is that block tags are not supported.

### Example:
```
Service {
    Name = "dev"
    Addr = ":8080"
    Insecure = true
    Deny = ["Reload", "Shutdown"]
}

Service {
    Name = "prod"
    Addr = ":80"
    ID = 49283

    Metrics {
        Route = "/metrics"
        Addr = ":8089"
    }
}
```