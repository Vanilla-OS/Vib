---
Title: Built-in modules
Description: Learn about the built-in modules that come with Vib and how to use them in your recipes.
PublicationDate: 2024-02-13
Listed: true
Authors:
  - mirkobrombin
Tags:
  - modules
---

> **Note**
> At the time of writing, Vib is in active development and the list of built-in modules may grow over time. This article covers the modules available in Vib v0.8.1.

Vib supports a variety of built-in modules that you can use to build your recipes. These modules are designed to automate common tasks, such as installing packages, building software, and running custom scripts.

Before proceeding, make sure to familiarize yourself with [how modules work](/vib/en/use-modules) since this article assumes you have a basic understanding of the module structure and how to use them in your recipes.

To keep this article concise, we'll cover only the fields that are specific to each module type, so `name`, `type` and `source` will be omitted if they don't have any specific fields.

## Summary

- [Package manager](#package-manager)
- [CMake](#cmake)
- [Dpkg-buildpackage](#dpkg-buildpackage)
- [Dpkg](#dpkg)
- [Go](#go)
- [Make](#make)
- [Meson](#meson)
- [Shell](#shell)
- [Flatpak](#flatpak)

## Package manager

This module allow to install packages using the package manager using the repositories configured in the image. You can change the package manager by changing the value of the `type` field. The following are currently supported:

- `apt`: Debian-based systems.
- `dnf`: Red Hat-based systems.

The following specific fields are available:

- `source`: Defines the source of the packages.

### Example

```yaml
- name: install-utils
  type: apt # or any other supported package manager
  source:
    packages:
      - curl
      - git
```

In the context of this module, this directive also supports the `packages` and `paths` fields. The `packages` field is a list of package names to install, while the `paths` field is a list of paths to `.inst` files containing package names each on a new line:

```yaml
- name: install-utils
  type: apt # or any other supported package manager
  source:
    paths:
      - "./utils.inst"
      - "./more-utils.inst"
```

where `utils.inst` and `more-utils.inst` follow the format:

```plaintext
curl
git
```

### Apt

> **Note**
> The following options requires Vib v.0.5.0 or later.

The `apt` module, has some additional fields under the `options` key:

- noRecommends: If set to `true`, the recommended packages will not be installed.
- installSuggestions: If set to `true`, the suggested packages will be installed.
- fixMissing: If set to `true`, the package manager will attempt to fix broken dependencies.
- fixBroken: If set to `true`, the package manager will attempt to fix broken packages.

```yaml
- name: install-utils
  type: apt
  source:
    packages:
      - curl
      - git
  options:
    noRecommends: true
    installSuggestions: true
    fixMissing: true
    fixBroken: true
```

> **Note**
> The above options if set to `false`, might still be overridden by the package manager's configuration.

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
  buildflags:
  - "-Dfoo=bar"
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

## Flatpak

The Flatpak module installs Flatpak packages using the `flatpak` command.

The following specific fields are available:

- `system`: If configured, the module will install the applications system-wide.
- `user`: If configured, the module will install the applications user-wide.

### Example

```yaml
- name: install-flatpak-app
  type: flatpak
  system:
    repourl: "https://flathub.org/repo/flathub.flatpakrepo"
    reponame: "flathub"
    install:
      - "org.gnome.Epiphany"
    remove:
      - "org.gnome.Epiphany"
  user:
    repourl: "https://flathub.org/repo/flathub.flatpakrepo"
    reponame: "flathub"
    install:
      - "org.gnome.Epiphany"
    remove:
      - "org.gnome.Epiphany"
```