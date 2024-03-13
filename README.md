# azure-b2c-sdk-for-go
![Tests](https://github.com/Schumann-IT/azure-b2c-sdk-for-go/actions/workflows/test.yml/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-80.4%25-brightgreen)

## Getting started

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

To add the latest version to your `go.mod` file, execute the following command.

```bash
go get github.com/Schumann-IT/azure-b2c-sdk-for-go
```

For more detailed usage examples, please checkout
- [go-ieftool](https://github.com/Schumann-IT/go-ieftool)
- [Terraform Provider azureadb2c](https://github.com/Schumann-IT/terraform-provider-azureadb2c)

## Usage

* create config file
```yaml
- name: test
  settings:
    SomeVariable: SomeTestContent 
- name: prod
  settings:
    SomeVariable: SomeProdContent 
```
* create some policies
```pre
src/
├─ local/
│  ├─ base.xml 
│  ├─ signupsignin.xml 
│  ├─ passwordreset.xml
├─ base.xml 
├─ extension.xml 

```

### Build policies

```go
package main

import (
	"log"
	
	"github.com/schumann-it/azure-b2c-sdk-for-go"
)

func main() {
    service, err := b2c.NewService("environments.yaml", "src", "build")
    if err != nil {
		log.Fatalf("failed to create service: %w", err)
    }
	env := "environment_name"
    err = service.BuildPolicies("environment_name")
    if err != nil {
		log.Fatalf("failed to build policies for environment %s: %w", env, err)
    }
}
```

### Deploy policies
```go
package main

import (
	"log"
	
	"github.com/schumann-it/azure-b2c-sdk-for-go"
)

func main() {
    service, err := b2c.NewService("environments.yaml", "src", "build")
    if err != nil {
		log.Fatalf("failed to create service: %w", err)
    }

	env := "test"
    err = service.DeployPolicies(env)
    if err != nil {
		log.Fatalf("failed to deploy policies for environment %s: %w", env, err)
    }
}
```
