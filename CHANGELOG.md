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

## [1.0.4] - 2025-09-12

### Fixed
- **AsyncConfigManager**: Manual message templates are now preserved during config refresh
- Previously, manually added message templates were lost when AsyncConfigManager refreshed configuration
- Manual templates now persist across automatic config reloads
