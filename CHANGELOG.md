# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

# [Unreleased]

### Added
- 

### Fixed
- 

### Features
- 

## [1.0.2] - 2025-09-12

### Changed
- **BREAKING**: Changed `Response.Data` field from `interface{}` to `map[string]any` for better type safety
- **BREAKING**: Changed `Response.Meta` field from `interface{}` to `map[string]any` for better type safety
- Updated all `map[string]interface{}` references to `map[string]any` throughout the codebase
- Updated documentation in README.md to reflect the new type definitions
- Updated all test files to use `map[string]any` instead of `map[string]interface{}`
- Updated example files to use the new type definitions

### Breaking Changes
- `Response.Data` field type changed from `interface{}` to `map[string]any`
- `Response.Meta` field type changed from `interface{}` to `map[string]any`
- All `map[string]interface{}` references changed to `map[string]any`
- This may require updates in client code that directly accesses these fields
