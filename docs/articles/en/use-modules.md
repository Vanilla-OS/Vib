---
Title: How to Use Vib Modules
Description: Add built-in, nested, and custom modules to a Vib recipe.
PublicationDate: 2026-08-24
Listed: true
Authors:
  - mirkobrombin
  - kbdharun
  - NN708
Tags:
  - modules
---

A module produces one or more build steps inside an image stage. Modules run in
the order listed by the recipe.

## Module structure

Every module has a unique `name` and a `type`:

```yaml
- name: install-build-tools
  type: apt
  sources:
    - packages:
        - build-essential
        - git
```

The other fields depend on the module type. See
[Built-in Modules](/vib/en/built-in-modules) for the plugin-specific fields.

## Sources

Modules can download or copy resources before their commands run:

```yaml
- name: install-vanilla-tools
  type: shell
  sources:
    - type: tar
      url: https://github.com/Vanilla-OS/vanilla-tools/releases/download/v1.0.1/vanilla-tools-amd64.tar.gz
      checksum: aef32f07820e0993e534e6bccfa1a6daae6c8c6f0543d3e073f4f121f2ef2e31
  commands:
    - install -Dm755 /sources/install-vanilla-tools/vanilla-tools/lpkg /usr/bin/lpkg
```

Vib 1.1 mounts a module's prepared sources at `/sources/MODULE_NAME`. Source
types include `git`, `tar`, `file`, and `local`.

A Git source accepts one revision strategy:

```yaml
- name: apx-gui
  type: meson
  sources:
    - type: git
      url: https://github.com/Vanilla-OS/apx-gui
      branch: main
```

Use `branch`, `tag`, or `commit`. A `commit` can be a commit hash or `latest`.
Use `checksum` for downloaded files and archives whenever the source publishes
a stable digest.

Limit a source to selected architectures with:

```yaml
only-arches:
  - amd64
  - arm64
```

## Nested modules

A module can contain child modules that prepare its build dependencies:

```yaml
- name: build-application
  type: go
  source:
    type: git
    url: https://example.com/application.git
    tag: v1.0.0
  modules:
    - name: install-go
      type: apt
      sources:
        - packages:
            - golang-go
```

## Included module files

Reuse local or remote module lists with `includes`:

```yaml
- name: shared-modules
  type: includes
  includes:
    - modules/common.yml
    - gh:organization/repository:main:modules/desktop.yml
```

Included files are expanded at their position in the module list.

## Custom plugins

Place project plugins in `plugins/` or shared plugins in
`/usr/share/vib/plugins/`. See [Making a Plugin](/vib/en/making-plugin) for the
plugin interface.
