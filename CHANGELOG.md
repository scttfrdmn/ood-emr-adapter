# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Bumped `github.com/scttfrdmn/substrate` 0.45.2 → 0.65.0 and regenerated go.sum. The recorded 0.45.2 checksum no longer matched the module the proxy serves (upstream re-tag), which broke `go test -tags=integration` with a go.sum SECURITY ERROR. Integration tests now build and pass.

### Added
- Initial scaffold — OOD compute adapter for Amazon EMR Serverless, translating Open OnDemand job lifecycle calls to the EMR Serverless API.
- CLI commands: `submit` (JSON job spec from stdin → EMR Serverless job run, prints `<application-id>/<job-run-id>`), `status <id>` (OOD-normalized status), `delete <id>` (cancel a job run), and `info <id>` (full `GetJobRun` response as JSON).
- Unit tests for status state mapping.
- Substrate integration tests for the EMR Serverless job run lifecycle.
- CI workflow with pinned action SHAs.
