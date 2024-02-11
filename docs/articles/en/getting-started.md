---
Title: Getting Started
Description: How to start using Vib to build your Container images.
PublicationDate: 2024-02-11
Authors:
  - mirkobrombin
---

Vib is a powerful tool that allows you to create container images using a YAML recipe that defines the necessary steps to build the image through the use of predefined or custom modules.

## Requirements

To use Vib, there are no specific requirements; you just need a Linux\* operating system (Mac OS and Windows will be supported in the future). Optionally, you can install a container engine to test and publish the images created to a registry.

\* Currently, Vib requires a Linux distribution with glibc.

### Supported Container Engines

- [Docker](https://www.docker.com/)
- [Podman](https://podman.io/)

Other container engines might work but have not been tested. If you have tested Vib with another container engine, please report it to our community.

## Installation

Vib is distributed as a single binary, so there's no need to install any runtime or dependencies. You can download the latest version of Vib from the [GitHub releases page](https://github.com/Vanilla-OS/Vib). Once downloaded, make the file executable and move it to a directory included in your PATH.

The following commands will allow you to download and install Vib:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/vib
chmod +x vib
mv vib ~/.local/bin
```

If wget is not installed, you can use curl:

```bash
curl -L https://github.com/Vanilla-OS/Vib/releases/latest/download/vib -o vib
chmod +x vib
mv vib ~/.local/bin
```

## Usage

To start using Vib, create a `vib.yaml` file in a new directory. This file will contain the recipe for your container image.

```bash
mkdir my-vib-project
cd my-vib-project
touch vib.yaml
```

Here's an example `vib.yaml` file:

```yaml
base: debian:sid-slim
labels:
  maintainer: My Awesome Team
args:
  DEBIAN_FRONTEND: noninteractive
runs:
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

In this example, we're creating a container image based on `debian:sid-slim` with some custom labels and environment variables. We're also installing a custom module that uses the default `go` module to clone a Git repository and install dependencies of the `golang` module via `apt`.

Once you've created the `vib.yaml` file, you can run the command:

```bash
vib build vib.yaml
```

to turn your recipe into a Containerfile. Use that file to build the container image with your container engine. To streamline the process, you can use the `compile` command to build the container image directly:

```bash
vib compile vib.yaml docker
```

changing `docker` with the container engine you have installed.

The generated `Containerfile` is compatible with any container engine that supports the OCI format, so you can use it with any container engine that supports this standard. Refer to your container engine's documentation for further information.

## Next Steps

Now that you've learned how to create a container image with Vib, you can start experimenting with predefined and custom modules to create more complex container images. Check out the [documentation](https://docs.vanillaos.org/collections/vib) for more information on all of Vib's features.

We recommend starting with the documentation on the [recommended structure of a Vib project](https://docs.vanillaos.org/collections/vib/project-structure) to understand how to best leverage Vib in your projects.
