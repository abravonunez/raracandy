# Contributing to raracandy

Thank you for your interest in contributing to raracandy!

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/abravonunez/raracandy.git
   cd raracandy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   make build
   ```

## Development Workflow

### Available Make Commands

- `make help` - Show all available commands
- `make build` - Build for current platform
- `make test` - Run tests with coverage
- `make lint` - Run linter
- `make fmt` - Format code
- `make clean` - Clean build artifacts

### Running Tests

```bash
make test
```

### Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Ensure `make lint` passes

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Maintenance tasks

**Examples:**
```
feat: add support for Pokemon Red/Blue
fix: correct checksum calculation for bag items
docs: update installation instructions
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Make your changes
4. Run tests and linting (`make check`)
5. Commit using conventional commits
6. Push to your fork
7. Open a Pull Request

## Release Process

Releases are automated using [Release Please](https://github.com/googleapis/release-please):

1. Merge PRs with conventional commit messages to `main`
2. Release Please creates a release PR automatically
3. Review and merge the release PR
4. GitHub Actions builds and publishes binaries via GoReleaser

## Questions?

Feel free to open an issue for any questions or concerns.
