---
Title: Built-in modules
Description: Learn about the built-in modules that come with Vib and how to use them in your recipes.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
---

> **Note**
> At the time of writing, Vib is in active development and the list of built-in modules may change. We are expanding support for the modules we offer and many more are planned for Vib v0.5.0.

Vib supports a variety of built-in modules that you can use to build your recipes. These modules are designed to automate common tasks, such as installing packages, building software, and running custom scripts.

Before proceeding, make sure to familiarize yourself with [how modules work](/articles/en/use-modules) since this article assumes you have a basic understanding of the module structure and how to use them in your recipes.

To keep this article concise, we'll cover only the fields that are specific to each module type, so `name`, `type` and `source` will be omitted if they don't have any specific fields.

## Summary

- [Apt](#apt)
- [CMake](#cmake)
- [Dpkg-buildpackage](#dpkg-buildpackage)
- [Dpkg](#dpkg)
- [Go](#go)
- [Make](#make)
- [Meson](#meson)
- [Shell](#shell)

## Apt

The Apt module allow to install packages using the APT package manager using the repositories configured in the image. It's suitable for Debian-based distributions.

The following specific fields are available:

- `source`: Defines the source of the packages.

### Example

```yaml
- name: install-utils
  type: apt
  source:
    packages:
      - curl
      - git
```

In the context of this module, this directive also supports the `packages` and `paths` fields. The `packages` field is a list of package names to install, while the `paths` field is a list of paths to files containing package names each on a new line:

```yaml
- name: install-utils
  type: apt
  source:
    paths:
      - "./utils"
      - "./more-utils"
```

where `utils` and `more-utils` follow the format:

```plaintext
curl
git
```

## CMake

The CMake module builds a project using the CMake build system. It's suitable for projects that use CMake as their build configuration tool.

The following specific fields are available:

- `buildFlags`: Additional flags to pass to the `cmake` command.

### Example

```yaml
- name: example-cmake-project
  type: cmake
  buildflags: "-DCMAKE_BUILD_TYPE=Release"
  source:
    url: "https://example.com/example-project.tar.gz"
    type: tar
```

## Dpkg-buildpackage

This module builds Debian packages from source using `dpkg-buildpackage` and installs the resulting `.deb` packages.

The following specific fields are available:

- `source`: source of the Debian package source code.

### Example

```yaml
- name: build-deb-package
  type: dpkg-buildpackage
  source:
    url: "https://example.com/package-source.tar.gz"
    type: tar
```

## Dpkg

The Dpkg module installs `.deb` packages directly using `dpkg` and resolves dependencies using `apt`.

The following specific fields are available:

- `source`: source of the `.deb` package(s) to install.

### Example

```yaml
- name: install-custom-deb
  type: dpkg
  source:
    paths:
      - "./packages/my-package.deb"
```

## Go

The Go module compiles Go projects, allowing for customization through build variables and flags.

The following specific fields are available:

- `buildFlags`: Flags for the `go build` command.

### Example

```yaml
- name: example-go-app
  type: go
  buildflags: "-v"
  source:
    url: "https://example.com/go-app-source.tar.gz"
    type: tar
```

## Make

The Make module automates the build process for projects that use GNU Make.

The following specific fields are available:

- `buildFlags`: Additional flags for the `make` command.

### Example

```yaml
- name: example-make-project
  type: make
  buildflags: "all"
  source:
    url: "https://example.com/make-project-source.tar.gz"
    type: tar
```

## Meson

This module is used for building projects configured with the Meson build system.

The following specific fields are available:

- `buildFlags`: Additional flags to pass to the `meson` command.

### Example

```yaml
- name: example-meson-project
  type: meson
  buildflags: "-Dfoo=bar"
  source:
    url: "https://example.com/meson-project-source.tar.gz"
    type: tar
```

## Shell

The Shell module executes arbitrary shell commands, offering the most flexibility for custom operations.

The following specific fields are available:

- `commands`: A list of shell commands to execute.

### Example

```yaml
- name: custom-setup
  type: shell
  commands:
    - "echo Hello, World!"
    - "apt update && apt install -y curl"
```
