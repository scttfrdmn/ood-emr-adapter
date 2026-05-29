# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial scaffold — OOD compute adapter for Amazon EMR Serverless, translating Open OnDemand job lifecycle calls to the EMR Serverless API.
- CLI commands: `submit` (JSON job spec from stdin → EMR Serverless job run, prints `<application-id>/<job-run-id>`), `status <id>` (OOD-normalized status), `delete <id>` (cancel a job run), and `info <id>` (full `GetJobRun` response as JSON).
- Unit tests for status state mapping.
- Substrate integration tests for the EMR Serverless job run lifecycle.
- CI workflow with pinned action SHAs.
