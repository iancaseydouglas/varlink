# Varlink

Varlink is a tool for managing Terraform variables and environments. It simplifies the process of setting up Terraform environment variables based on .tfvars files in your project structure.

## Features

- Automatically finds and parses .tfvars files
- Sets Terraform environment variables (TF_VAR_*)
- Supports multiple environments (dev, staging, prod, etc.)
- Allows specifying search depth for .tfvars files
- Provides a dry-run mode to preview actions
- Can deactivate all Terraform-related environment variables

## Usage

```
varlink [flags]

Flags:
  -levels int
        Levels above 'environments' to search for tfvars (default 1)
  -max-depth int
        Maximum directory depth to search (default 10)
  -deactivate
        Deactivate Terraform environment variables
  -dry-run
        Show what would be done without making changes
```
## Assumed directory structure

├── environments
│   ├── dev
│   │   ├── dev.tfvars
│   │   ├── service1
│   │   │   ├── main.tf
│   │   │   └── service1.tfvars
│   │   └── service2
│   │       ├── main.tf
│   │       └── service2.tfvars
│   ├── staging
│   │   ├── main.tf
│   │   └── staging.tfvars
│   └── prod
│       ├── main.tf
│       └── prod.tfvars
├── global.tfvars

## Building

TODO

```
```

## Testing

TODO


```
```


```
```


