---
Title: Getting Started
Description: Install Vib 1.1 and build a container image from a recipe.
PublicationDate: 2026-08-24
Listed: true
Authors:
  - mirkobrombin
  - kbdharun
  - surinameclubcard
  - NN708
Tags:
  - getting-started
---

Vib generates container images from YAML recipes. A recipe defines image stages
and the modules that configure each stage.

## Requirements

Vib runs on a glibc-based Linux distribution. Install Docker or Podman only if
you want Vib to build the generated Containerfile for you.

## Install Vib

Download the binary for your architecture from the
[Vib releases page](https://github.com/Vanilla-OS/Vib/releases). For AMD64:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/vib-amd64
chmod +x vib-amd64
mkdir -p ~/.local/bin
mv vib-amd64 ~/.local/bin/vib
```

Download the matching plugin archive and install it globally:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/plugins-amd64.tar.gz
sudo mkdir -p /usr/share/vib/plugins
sudo tar -xvf plugins-amd64.tar.gz -C /usr/share/vib/plugins --strip-components=2
```

Use the `arm64` assets on an AArch64 system. Project-specific plugins can be
placed in a `plugins` directory beside the recipe instead.

## Create a recipe

Create `recipe.yml`:

```yaml
name: My Node Image
id: my-node-image
vibversion: 1.1.0

stages:
  - id: build
    base: node:current-slim
    labels:
      maintainer: My Team
    expose:
      "3000": "tcp"
    entrypoint:
      exec:
        - node
        - /app/app.js
    modules:
      - name: build-app
        type: shell
        sources:
          - type: git
            url: https://github.com/mirkobrombin/node-sample
            branch: main
        commands:
          - mv /sources/build-app/node-sample /app
          - cd /app
          - npm install
          - npm run build
```

Vib 1.1 mounts module sources under `/sources/MODULE_NAME` while generating the
corresponding build step.

## Validate and build

Check that Vib can load the recipe, then generate a Containerfile:

```bash
vib test recipe.yml
vib build recipe.yml
```

Choose a different output name with `--output`, or `-o`:

```bash
vib build --output Containerfile.reunion recipe.yml
```

Generate the Containerfile and build the image in one command:

```bash
vib compile --runtime docker recipe.yml
```

Docker and Podman are supported. If `--runtime` is omitted, Vib uses Docker
when available, then Podman. On Vanilla OS, run `vib compile` inside
`host-shell` so it can access the host container engine.

## Next steps

- [Recipe Structure](/vib/en/recipe-structure)
- [Project Structure](/vib/en/project-structure)
- [How to Use Vib Modules](/vib/en/use-modules)
- [Built-in Modules](/vib/en/built-in-modules)
