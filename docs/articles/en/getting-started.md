---
Title: Getting Started
Description: How to start using Vib to build your Container images.
PublicationDate: 2024-02-11
Listed: true
Authors:
  - mirkobrombin
  - kbdharun
Tags:
  - getting-started
---

Vib is a powerful tool that allows you to create container images using a YAML recipe that defines the necessary steps to build the image through the use of predefined or custom modules.

## Requirements

To use Vib, there are no specific requirements; you just need a Linux\* operating system (Mac OS and Windows will be supported in the future). Optionally, you can install a container engine to test and publish the images created to a registry.

\* Currently, Vib requires a Linux distribution with `glibc`.

### Supported Container Engines

- [Docker](https://www.docker.com/)
- [Podman](https://podman.io/)

Other container engines might work but have not been tested. If you have tested Vib with another container engine, please report it to our community.

## Installation

Vib is distributed as a single binary, so there's no need to install any runtime or dependencies. You can download the latest version of Vib from the [GitHub releases page](https://github.com/Vanilla-OS/Vib). In addition to this, Vib has official plugins which are used for all the Vanilla-OS images, they can also be downlaoded from the [Github releases page](https://github.com/Vanilla-OS/Vib) as the `plugins.tar.xz` archvie. Once downloaded, make vib executable and move it to a directory included in your PATH. Vib searches for plugins in a global search path at `/usr/share/vib/plugins/` and inside the `plugins` directory in your project directory. It is recommended to extract `plugins.tar.xz` to `/usr/share/vib/plugins/` as they are considered core vib plugins and may be used by a lot of images.

The following commands will allow you to download and install Vib:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/vib
chmod +x vib
mv vib ~/.local/bin
```

If wget is not installed, you can use curl:

```bash
curl -SLO https://github.com/Vanilla-OS/Vib/releases/latest/download/vib
chmod +x vib
mv vib ~/.local/bin
```

The following commands for the plugins:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/plugins.tar.xz
mkdir -p /usr/share/vib/plugins
tar -xvf plugins.tar.xz -C /usr/share/vib/plugins/
```

Or with curl:

```bash
curl -SLO https://github.com/Vanilla-OS/Vib/releases/latest/download/plugins.tar.xz
mkdir -p /usr/share/vib/plugins
tar xvf plugins.tar.xz -C /usr/share/vib/plugins
```

## Usage

To start using Vib, create a `vib.yml` file in a new directory. This file will contain the recipe for your container image.

```bash
mkdir my-vib-project
cd my-vib-project
touch vib.yml
```

Here's an example `vib.yml` file:

```yml
name: My Image
id: my-image
stages:
  - id: build
    base: debian:sid-slim
    singlelayer: false
    labels:
      maintainer: My Awesome Team
    args:
      DEBIAN_FRONTEND: noninteractive
    runs:
      commands:
        - echo 'APT::Install-Recommends "0";' > /etc/apt/apt.conf.d/01norecommends
    modules:
      - name: update
        type: shell
        commands:
          - apt update
      - name: vib
        type: go
        source:
          type: git
          url: https://github.com/vanilla-os/vib
          branch: main
          commit: latest
        buildVars:
          GO_OUTPUT_BIN: /usr/bin/vib
        modules:
          - name: golang
            type: apt
            source:
              packages:
                - golang
                - ca-certificates
```

In this example, we're creating a container image with one stage based on `debian:sid-slim` with some custom labels and environment variables. We're also installing a custom module that uses the default `go` module to clone a Git repository and install dependencies of the `golang` module via `apt`.

Once you've created the `vib.yaml` file, you can run the command:

```bash
vib build vib.yaml
```

to turn your recipe into a Containerfile. Use that file to build the container image with your container engine. To streamline the process, you can use the `compile` command to build the container image directly:

```bash
vib compile vib.yml --runtime docker
```

changing `docker` with the container engine you have installed. Both `docker` and `podman` are supported. If you leave out the `--runtime` flag, Vib will use the default container engine giving priority to Docker.

> **Note:**
> For versions of Vib before 0.5.0, the syntax of the `compile` command was different. The `--runtime` flag was not available, and the command was `vib compile vib.yml docker`.

The generated `Containerfile` is compatible with both Docker and Podman.

## Next Steps

Now that you've learned how to create a container image with Vib, you can start experimenting with predefined and custom modules to create more complex container images. Check out the [documentation](/collections/vib) for more information on all of Vib's features.

We recommend starting with the documentation on the [recommended structure of a Vib project](/vib/en/project-structure) to understand how to best leverage Vib in your projects.
