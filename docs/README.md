# Swagger/OpenAPI Documentation

This directory contains the generated Swagger (OpenAPI) documentation for the users-adapter-gin module.

- `swagger.yaml` and `swagger.json` are the OpenAPI specs.
- Docs are generated from code annotations using [swaggo/swag](https://github.com/swaggo/swag).

## How to update docs

1. Install swag CLI: `go install github.com/swaggo/swag/cmd/swag@latest`
2. Run `swag init` in the root of this module to regenerate docs from code comments.
3. Serve Swagger UI at `/swagger/*any` in your Gin app (see `main.go`).
