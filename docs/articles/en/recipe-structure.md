---
Title: Structure of a Vib Recipe
Description: Reference for Vib 1.1 recipe and stage fields.
PublicationDate: 2026-08-24
Listed: true
Authors:
  - mirkobrombin
  - kbdharun
  - lambdaclan
  - NN708
Tags:
  - modules
  - recipe
---

A Vib recipe contains image metadata, one or more stages, and optional finalize
plugins.

```yaml
name: Example Image
id: example-image
vibversion: 1.1.0
includespath: includes.container

stages:
  - id: build
    base: golang:bookworm
    args:
      DEBIAN_FRONTEND: noninteractive
    modules:
      - name: build-app
        type: shell
        commands:
          - go build -o /output/app ./cmd/app

  - id: runtime
    base: debian:bookworm-slim
    labels:
      org.opencontainers.image.title: Example Image
    copy:
      - from: build
        srcdst:
          /output/app: /usr/bin/app
    cleanup:
      - /var/lib/apt/lists/*
    expose:
      "8080": "tcp"
    entrypoint:
      exec:
        - /usr/bin/app
```

## Recipe fields

- `name`: display name of the image.
- `id`: stable image identifier.
- `vibversion`: recipe format version. Vib 1.1 accepts versions starting at 1.0.0.
- `stages`: ordered image stages.
- `includespath`: optional path used in place of `includes.container`.
- `finalize`: optional plugins run after the image build.

## Stage fields

- `id`: unique stage identifier.
- `base`: source image, including `scratch` where appropriate.
- `labels`: OCI image labels.
- `env`: environment values written to the image.
- `args`: build arguments.
- `runs`: commands run before modules.
- `modules`: ordered build modules.
- `adds`: host files or directories added to the stage.
- `addincludes`: add the `includespath` filesystem tree to the stage.
- `copy`: files copied from the host or another stage.
- `cleanup`: paths removed after generated module and run commands.
- `expose`: port and protocol mapping.
- `cmd`: default container command.
- `entrypoint`: container entry point.

## Working directories

`runs`, `modules`, `adds`, `copy`, `cmd`, and `entrypoint` accept `workdir` in
their respective structures. For example:

```yaml
runs:
  workdir: /app
  commands:
    - ./configure

modules:
  - name: build
    type: shell
    workdir: /app
    commands:
      - make
```

## Add host files

Map project files to image paths with `adds`:

```yaml
adds:
  - srcdst:
      files/example.conf: /etc/example.conf
```

For a complete filesystem tree, put files under `includes.container/` and set
`addincludes: true` on the target stage. See
[Project Structure](/vib/en/project-structure).

## Copy between stages

Use the source stage ID in `from`:

```yaml
copy:
  - from: build
    srcdst:
      /output/app: /usr/bin/app
```

## Commands

`cmd` and `entrypoint` use an `exec` list and an optional `workdir`:

```yaml
entrypoint:
  workdir: /app
  exec:
    - /usr/bin/app
    - --serve
```

See [How to Use Vib Modules](/vib/en/use-modules) for module fields.
