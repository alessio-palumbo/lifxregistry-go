# lifxregistry-go

This repository provides a Go-based generator for the [LIFX product registry](https://github.com/LIFX/products), producing Go types and definitions for use in applications that discover or interact with known LIFX device models.

It includes:

- `cmd/registry-gen` — the generator tool for the LIFX product registry.
- `cmd/registry-gen/src/registry.json` — the upstream products registry definition from the official [LIFX products repo](https://github.com/LIFX/products).
- `cmd/registry-gen/src/registry_commit.txt` — the upstream commit the current definition was sourced from.
- `cmd/registry-gen/gen/` — generated Go files containing typed device metadata.

## Usage

To update and regenerate the LIFX registry Go code:

```bash
make regenerate
```

This will:

- Clone the latest upstream LIFX/products repo.
- Check if the products commit has changed.
- If changed, copy the new products.json and commit hash.
- Regenerate Go code under gen/.

### Requirements

- Go 1.24+
- Git
- Internet connection (to fetch latest upstream repo)

### License

This project is MIT licensed. See LICENSE for details.
