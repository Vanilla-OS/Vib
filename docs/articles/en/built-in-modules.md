---
Title: Built-in Modules
Description: Reference for the modules distributed with Vib 1.1.
PublicationDate: 2026-08-24
Listed: true
Authors:
  - mirkobrombin
  - NN708
Tags:
  - modules
---

Vib 1.1 distributes build plugins for APT, CMake, dpkg-buildpackage, Flatpak,
Go, Make, and Meson. The `shell` and `includes` modules are part of Vib itself.

Every module needs a unique `name` and a `type`. Sources used by a module are
mounted under `/sources/MODULE_NAME` while its generated build command runs.

## Source fields

A source can define:

- `type`: `git`, `tar`, `file`, or `local`.
- `url`: remote source URL.
- `checksum`: SHA-256 checksum for a file or archive.
- `branch`, `tag`, or `commit`: Git revision selector.
- `packages`: package names for a package-manager plugin.
- `path`: local source path or package-list file.
- `only-arches`: CPU architectures allowed to use the source.

## APT

Install Debian packages with a list or a local `.inst` file:

```yaml
- name: install-utils
  type: apt
  sources:
    - packages:
        - curl
        - git
    - path: ./extra-packages.inst
  options:
    no_recommends: true
    install_suggests: false
    fix_missing: false
    fix_broken: false
```

Each line in an `.inst` file is passed as a package name. The supported option
keys map to the corresponding `apt-get install` flags.

## CMake

```yaml
- name: example-cmake-project
  type: cmake
  buildflags: -DCMAKE_BUILD_TYPE=Release
  source:
    type: tar
    url: https://example.com/example-project.tar.gz
    checksum: SHA256
```

`buildvars` can provide build variables used by the plugin.

## dpkg-buildpackage

Build a Debian source package and install its resulting packages:

```yaml
- name: example-debian-package
  type: dpkg-buildpackage
  source:
    type: git
    url: https://example.com/example-debian-package.git
    tag: v1.0.0
```

The source must contain valid Debian packaging metadata.

## Go

```yaml
- name: example-go-app
  type: go
  buildflags: -trimpath
  buildvars:
    GO_OUTPUT_BIN: /usr/bin/example-go-app
  source:
    type: git
    url: https://example.com/example-go-app.git
    commit: latest
```

`GO_OUTPUT_BIN` sets the output path. Without it, the module name is used.

## Make

```yaml
- name: example-make-project
  type: make
  buildcommand: make PREFIX=/usr
  intermediatesteps:
    - make test
  installcommand: make PREFIX=/usr install
  sources:
    - type: tar
      url: https://example.com/example-make-project.tar.gz
      checksum: SHA256
```

The defaults are `make` and `make install`.

## Meson

```yaml
- name: example-meson-project
  type: meson
  buildflags:
    - -Dfoo=enabled
  sources:
    - type: tar
      url: https://example.com/example-meson-project.tar.gz
      checksum: SHA256
```

## Shell

Run custom commands in order:

```yaml
- name: custom-setup
  type: shell
  sources:
    - type: file
      url: https://example.com/example.conf
      checksum: SHA256
  commands:
    - install -Dm644 /sources/custom-setup/example.conf /etc/example.conf
  cleanup:
    - /tmp/example-cache
```

Module-level `cleanup` paths are removed after the module commands.

## Flatpak

Configure system or user Flatpak remotes and application lists:

```yaml
- name: install-flatpak-apps
  type: flatpak
  system:
    repo-url: https://flathub.org/repo/flathub.flatpakrepo
    repo-name: flathub
    install:
      - org.gnome.Epiphany
    remove: []
```

The `user` block accepts the same fields. The plugin creates setup services so
the selected applications are installed when the image runs.

## Includes

Insert modules stored in local or remote YAML files:

```yaml
- name: shared-modules
  type: includes
  includes:
    - modules/common.yml
    - gh:organization/repository:main:modules/desktop.yml
```

Included modules run in the listed order.
