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

## [1.0.1] - 2025-09-12

### Changed
- Renamed `responseBuilder` struct to `ResponseBuilder` for proper Go naming conventions
- Updated all method signatures and return types to use `ResponseBuilder`
- Updated documentation to reflect the new naming convention
- Fixed linting warning in test file by using `context.TODO()` instead of `nil`

### Breaking Changes
- `responseBuilder` struct is now exported as `ResponseBuilder`
- All method signatures now return `*ResponseBuilder` instead of `*responseBuilder`
- `BuildResponse()` method now accepts `*ResponseBuilder` instead of `*responseBuilder`
- `ParseResponseBuilderError()` now returns `(*ResponseBuilder, bool)` instead of `(*responseBuilder, bool)`
