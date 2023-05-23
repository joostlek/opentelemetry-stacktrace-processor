# OpenTelemetry stacktrace processor

[![GitHub Release][releases-shield]][releases]
[![Go Reference](https://pkg.go.dev/badge/github.com/joostlek/opentelemetry-stacktrace-processor.svg)](https://pkg.go.dev/github.com/joostlek/opentelemetry-stacktrace-processor)
![Project Stage][project-stage-shield]
![Project Maintenance][maintenance-shield]
[![License][license-shield]](LICENSE.md)

[![Build Status][build-shield]][build]
[![Code Coverage][codecov-shield]][codecov]
[![Code Smells][code-smells]][sonarcloud]

Transform stacktraces from javascript into readable stacktraces to improve OpenTelemetry usability.

## About

When collecting stacktraces from javascript based application, the stacktraces are usually based on minified or compiled code.

Example:
```
padStart@http://localhost:4203/lineSlicer.min.js:1:228
padStart@http://localhost:4203/cubeSlicer.min.js:1:228
```

To be able to know where the actual problem is, you should combine this with a source map.
Commercial tracing solutions provide ways to upload source maps to their dashboard, and they figure it out.

This processor is made so the generic OpenTelemetry collector can have the same functionality.
Place the source maps in a directory and the processor will read them at startup.
The end result:

```
padStart@lib/lineSlicer.js:11:53
padStart@http://localhost:4203/cubeSlicer.min.js:1:228
```
In this example the map for lineSlicer.js is available, but not for cubeSlicer.
In that case the line is ignored, so there's always a full picture of the flow present.

## Installation

Install as [described in the OpenTelemetry docs](https://opentelemetry.io/docs/collector/custom-collector/)

## Changelog & Releases

This repository keeps a change log using [GitHub's releases][releases]
functionality. The format of the log is based on
[Keep a Changelog][keepchangelog].

Releases are based on [Semantic Versioning][semver], and use the format
of ``MAJOR.MINOR.PATCH``. In a nutshell, the version will be incremented
based on the following:

- ``MAJOR``: Incompatible or major changes.
- ``MINOR``: Backwards-compatible new features and enhancements.
- ``PATCH``: Backwards-compatible bugfixes and package updates.

## Contributing

This is an active open-source project. We are always open to people who want to
use the code or contribute to it.

We've set up a separate document for our
[contribution guidelines](.github/CONTRIBUTING.md).

Thank you for being involved! :heart_eyes:

## Setting up development environment

This Go project relies on Go. For development improvements, Python is used.

You need at least:

- Go 1.19
- Python 3.10+
- [Poetry][poetry-install]

```bash
poetry install
```

As this repository uses the [pre-commit][pre-commit] framework, all changes
are linted and tested with each commit. You can run all checks and tests
manually, using the following command:

```bash
poetry run pre-commit run --all-files
```

## Authors & contributors

The original setup of this repository is by [Joost Lekkerkerker][joostlek]

For a full list of all authors and contributors,
check [the contributor's page][contributors].

## License

MIT License

Copyright (c) 2023 Joost Lekkerkerker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

[build-shield]: https://github.com/joostlek/opentelemetry-stacktrace-processor/actions/workflows/tests.yaml/badge.svg
[build]: https://github.com/joostlek/opentelemetry-stacktrace-processor/actions
[code-smells]: https://sonarcloud.io/api/project_badges/measure?project=joostlek_opentelemetry-stacktrace-processor&metric=code_smells
[codecov-shield]: https://codecov.io/gh/joostlek/opentelemetry-stacktrace-processor/branch/master/graph/badge.svg
[codecov]: https://codecov.io/gh/joostlek/opentelemetry-stacktrace-processor
[commits-shield]: https://img.shields.io/github/commit-activity/y/joostlek/opentelemetry-stacktrace-processor.svg
[commits]: https://github.com/joostlek/opentelemetry-stacktrace-processor/commits/master
[contributors]: https://github.com/joostlek/opentelemetry-stacktrace-processor/graphs/contributors
[joostlek]: https://github.com/joostlek
[keepchangelog]: http://keepachangelog.com/en/1.0.0/
[license-shield]: https://img.shields.io/github/license/joostlek/opentelemetry-stacktrace-processor.svg
[maintenance-shield]: https://img.shields.io/maintenance/yes/2023.svg
[poetry-install]: https://python-poetry.org/docs/#installation
[poetry]: https://python-poetry.org
[pre-commit]: https://pre-commit.com/
[project-stage-shield]: https://img.shields.io/badge/project%20stage-experimental-yellow.svg
[python-versions-shield]: https://img.shields.io/pypi/pyversions/opentelemetry-stacktrace-processor
[releases-shield]: https://img.shields.io/github/release/joostlek/opentelemetry-stacktrace-processor.svg
[releases]: https://github.com/joostlek/opentelemetry-stacktrace-processor/releases
[semver]: http://semver.org/spec/v2.0.0.html
[sonarcloud]: https://sonarcloud.io/summary/new_code?id=joostlek_opentelemetry-stacktrace-processor
