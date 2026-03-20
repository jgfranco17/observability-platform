# Observability Platform Backend

![STATUS](https://img.shields.io/badge/status-active-brightgreen?style=for-the-badge)
![LICENSE](https://img.shields.io/badge/license-BSD3-blue?style=for-the-badge)

---

## Introduction

A simple observability platform backend implemented in Go. This service provides a REST
API for storing and retrieving observability reports, designed for easy integration with
various monitoring tools.

## Deployment

> [!NOTE]
> Deployment site pending.

## Testing

```bash
# Run standard assertions with go-test
just test
```

### Automation

#### GitHub Actions Integration

Tests are automatically run on:

- Every pull request
- Every push to main branch
- Scheduled nightly runs

#### Quality Gates

- All tests must pass before merging
- Minimum code coverage requirements
- Performance benchmarks must be met
