# WIP - Config Format

> This package implements an unmarshaler for a simplified version of the HCL configuration language.

### Goals & Differences to HCL.

* No tag blocks.
* No variables.
* Support for `encoding.TextMarshaler` & `encoding.TextUnmarshaler`.
* Improved error messages.

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

``` go
type Metrics struct {
	Route string
	Addr  string
}

type Service struct {
	Name    string
	Addr    string
	Deny    []string
	Metrics *Metrics
}

type Config struct {
	Service []*Service
}

func main() {
	var c Config
	data, _ = iotuil.ReadAll("services.conf")
	_ = config.Unmarshal(data, &c)
}

```
