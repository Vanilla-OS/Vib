---
Title: Structure of a Vib recipe
Description: Learn about the structure of a Vib recipe.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
Tags:
  - modules
  - recipe
---

A Vib recipe is a YAML file that contains the instructions to build a container image. It's composed of three blocks:

- metadata
- configuration
- modules

The following is a complete example of a Vib recipe:

```yaml
# metadata
base: debian:sid-slim
name: My Image
id: my-image-id

# configuration
singlelayer: false
labels:
  maintainer: My Awesome Team
adds:
  - /path/to/add
args:
  - arg1: value1
  - arg2: value2
runs:
  - echo "Hello, World!"
expose: 8080
cmd: /bin/bash
entrypoint:
  - /bin/bash

# modules
modules:
  - name: update
    type: shell
    commands:
      - apt update
```

## Metadata

The metadata block contains the following mandatory fields:

- `base`: the base image to start from, can be any Docker image from any registry or even `scratch`.
- `name`: the name of the image.

The following fields are optional:

- `id`: the ID of the image, can be used by platforms like [Atlas](https://images.vanillaos.org/#/) to identify the image.

## Configuration

The configuration block contains the following optional fields:

- `singlelayer`: a boolean value that indicates if the image should be built as a single layer. This is useful in some cases to reduce the size of the image (e.g. when building an image using a rootfs, an example [here](https://github.com/Vanilla-OS/pico-image/blob/5b0e064677f78f6e89d619dcb4df4e585bef378f/recipe.yml)).
- `labels`: a map of labels to apply to the image, useful to add metadata to the image that can be read by the container runtime.
- `adds`: a list of files or directories to add to the image, useful to include files in the image that are not part of the source code (the preferred way to include files in the image is to use the `includes.container/` directory, see [Project Structure](/docs/articles/en/project-structure)).
- `args`: a list of environment variables to set in the image.
- `runs`: a list of commands to run in the image (as an alternative to the `shell` module, useful for dividing the commands of your recipe from those needed to configure the image, for example to disable the recommended packages in apt).
- `expose`: a list of ports to expose in the image.
- `cmd`: the command to run when the container starts.
- `entrypoint`: the entry point for the container, it's similar to `cmd` but it's not overridden by the command passed to the container at runtime, useful to handle the container as an executable.

## Modules

The modules block contains a list of modules to use in the recipe. Each module is a YAML snippet that defines a set of instructions. The common structure is:

```yaml
- name: name-of-the-module
  type: type-of-the-module
  # specific fields for the module type
```

Refer to the [Use Modules](/vib/en/use-modules) article for more information on how to use modules in a recipe and [Built-in Modules](/vib/en/built-in-modules) for a list of the built-in modules and their specific fields.

You can also write your own modules by making a Vib plugin, see the [Making a Plugin](/vib/en/making-plugin) article for more information.
