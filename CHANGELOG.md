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

## [1.0.5] - 2025-09-18

### Added
- **Context Helper Functions**: Added utility functions for extracting language and protocol from context
  - `GetLanguageFromContext(ctx context.Context) (string, bool)` - Extract language with found indicator
  - `GetProtocolFromContext(ctx context.Context) (string, bool)` - Extract protocol with found indicator
  - `GetLanguage(ctx context.Context) string` - Extract language directly (returns empty string if not found)
  - `GetProtocol(ctx context.Context) string` - Extract protocol directly (returns empty string if not found)
  - All functions handle nil context gracefully and include comprehensive test coverage
