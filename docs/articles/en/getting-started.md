---
Title: Getting Started
Description: How to start using Vib to build your Container images.
PublicationDate: 2024-12-28
Listed: true
Authors:
  - mirkobrombin
  - kbdharun
  - surinameclubcard
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

Vib is distributed as a single binary, so there's no need to install any runtime or dependencies. You can download the latest version of Vib from the [GitHub releases page](https://github.com/Vanilla-OS/Vib/releases). In addition to this, Vib has official plugins which are used for all the Vanilla-OS images, they can also be downlaoded from the [Github releases page](https://github.com/Vanilla-OS/Vib/releases) as the `plugins-*.tar.gz` archive. Once downloaded, make `vib` executable and move it to a directory included in your `PATH`. Vib searches for plugins in a global search path at `/usr/share/vib/plugins/` and inside the `plugins` directory in your project directory. It is recommended to extract `plugins-*.tar.gz` to `/usr/share/vib/plugins/` as they are considered core vib plugins and may be used by a lot of images.

The following commands will allow you to download and install Vib (_supported architectures amd64, arm64_):

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/vib-amd64
chmod +x vib-amd64
mv vib-amd64 ~/.local/bin/vib
```

If wget is not installed, you can use curl:

```bash
curl -SLO https://github.com/Vanilla-OS/Vib/releases/latest/download/vib-amd64
chmod +x vib-amd64
mv vib-amd64 ~/.local/bin/vib
```

The following commands for the plugins:

```bash
wget https://github.com/Vanilla-OS/Vib/releases/latest/download/plugins-amd64.tar.gz
sudo mkdir -p /usr/share/vib/plugins
sudo tar -xvf plugins-amd64.tar.gz -C /usr/share/vib/plugins/ --strip-components=2
```

Or with curl:

```bash
curl -SLO https://github.com/Vanilla-OS/Vib/releases/latest/download/plugins-amd64.tar.gz
sudo mkdir -p /usr/share/vib/plugins
sudo tar -xvf plugins-amd64.tar.gz -C /usr/share/vib/plugins --strip-components=2
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
name: my-recipe
id: my-node-app
stages:
  - id: build
    base: node:current-slim
    labels:
      maintainer: My Awesome Team
    args:
      DEBIAN_FRONTEND: noninteractive
    expose: 
      "3000": ""
    entrypoint:
      exec:
        - node
        - /app/app.js
    runs:
      commands:
        - echo 'APT::Install-Recommends "0";' > /etc/apt/apt.conf.d/01norecommends
    modules:
    - name: build-app
      type: shell
      source:
        type: git
        url: https://github.com/mirkobrombin/node-sample
        branch: main
        commit: latest
      commands:
        - mv /sources/build-app /app
        - cd /app
        - npm i
        - npm run build
```

In this example, we're creating a container image with one stage based on `node:current-slim` with some custom labels and environment variables. The image uses a single module to build a Node.js application from a Git repository. The application is then installed and built using `npm`. Then it exposes the port `3000` and sets the entrypoint to node `/app/app.js`. 

Once you've created the `vib.yml` file, you can run the command:

```bash
vib build vib.yml
```

to turn your recipe into a Containerfile. Use that file to build the container image with your container engine. To streamline the process, you can use the `compile` command to build the container image directly:

```bash
vib compile --runtime docker
```

changing `docker` with the container engine you have installed. Both `docker` and `podman` are supported. If you leave out the `--runtime` flag, Vib will use the default container engine giving priority to Docker.

> **Note:**
> On a Vanilla OS host, you need to run `vib compile` from the `host-shell`.

The generated `Containerfile` is compatible with both Docker and Podman.

## Next Steps

Now that you've learned how to create a container image with Vib, you can start experimenting with predefined and custom modules to create more complex container images. Check out the [documentation](/collections/vib) for more information on all of Vib's features.

We recommend starting with the documentation on the [recommended structure of a Vib project](/vib/en/project-structure) to understand how to best leverage Vib in your projects.
