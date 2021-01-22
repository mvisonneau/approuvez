# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org).

## [Unreleased]

## [v0.1.0] - 2021-01-22

### Added

- gRPC server/client based architecture over mTLS
- documentation!

### Changed/Removed

- Refactored and dropped the support/architecture based upon lambda/apigateway/websockets
- Reduced the featureset to release a more robust iteration of the app
- Updated all dependencies

## [v0.0.1] - 2020-12-23

### Added

- Command line client binary to use as part of CI jobs
- Lambda function to handle the slack callbacks
- Lambda function to handle websocket interactions with the clients
- End-to-end Terraform configuration for the AWS components

[Unreleased]: https://github.com/mvisonneau/approuvez/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/mvisonneau/approuvez/tree/v0.1.0
[v0.0.1]: https://github.com/mvisonneau/approuvez/tree/v0.0.1