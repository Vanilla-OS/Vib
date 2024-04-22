---
Title: Structure of a Vib recipe
Description: Learn about the structure of a Vib recipe.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
  - kbdharun
Tags:
  - modules
  - recipe
---

> **Note**
> Stages were introduced in Vib v0.6.0, if you are using an older version, please keep in mind all the stage fields are at the top level of the recipe, so no multiple stages are supported.

A Vib recipe is a YAML file that contains the instructions to build a container image. It's composed of two blocks:

- metadata
- stages

The following is a complete example of a Vib recipe:

```yml
# metadata
name: My Image
id: my-image-id

# stages
stages:
  - id: build
    base: debian:sid-slim
    singlelayer: false
    labels:
      maintainer: My Awesome Team
    adds:
      - /extra/path/to/add
    args:
      - arg1: value1
      - arg2: value2
    runs:
      - some-random-command --that-must-run --on-top-of-all modules
    modules:
      - name: build
        type: go
        buildvars:
          GO_OUTPUT_BIN: "/path/to/output"
        source:
          url: https://github.com/my-awesome-team/my-awesome-repo
          type: git
          branch: main
          commit: sdb997f0eeb67deaa5940f7c31a19fe1101d3d49
        modules:
        - name: build-deps
          type: apt
          source:
            packages:
            - golang-go

  - id: dist
    base: debian:sid-slim
    singlelayer: false
    labels:
      maintainer: My Awesome Team
    expose: 
      "8080": "tcp"
      "8081": ""
    entrypoint: ["/app"]
    copy:
      - from: build
        src: /path/to/output
        dest: /app
    cmd: ["/app"]
    modules:
      - name: run
        type: shell
        commands:
          - ls -la /app
```

## Metadata

The metadata block contains the following mandatory fields:

- `base`: the base image to start from, can be any Docker image from any registry or even `scratch`.
- `name`: the name of the image.
- `id`: the ID of the image is used to specify an image's unique identifier, it is used by platforms like [Atlas](https://images.vanillaos.org/#/) to identify the image.
- `stages`: a list of stages to build the image, useful to split the build process into multiple stages (e.g. to build the application in one stage and copy the artifacts into another one).

## Stages

Stages are a list of instructions to build an image, useful to split the build process into multiple stages (e.g. to build the application in one stage and copy the artifacts into another one). Each stage is a YAML snippet that defines a set of instructions.

Each stage has the following fields:

- `singlelayer`: a boolean value that indicates if the image should be built as a single layer. This is useful in some cases to reduce the size of the image (e.g. when building an image using a rootfs, an example [here](https://github.com/Vanilla-OS/pico-image/blob/5b0e064677f78f6e89d619dcb4df4e585bef378f/recipe.yml)).
- `labels`: a map of labels to apply to the image, useful to add metadata to the image that can be read by the container runtime.
- `adds`: a list of files or directories to add to the image, useful to include files in the image that are not part of the source code (the preferred way to include files in the image is to use the `includes.container/` directory, see [Project Structure](/docs/articles/en/project-structure)).
- `args`: a list of environment variables to set in the image.
- `runs`: a list of commands to run in the image (as an alternative to the `shell` module, useful for dividing the commands of your recipe from those needed to configure the image, for example, to disable the recommended packages in apt).
- `expose`: a list of ports to expose in the image.
- `cmd`: the command to run when the container starts.
- `entrypoint`: the entry point for the container, it's similar to `cmd` but it's not overridden by the command passed to the container at runtime, useful to handle the container as an executable.
- `copy`: a list of files or directories to copy from another stage, useful to copy files from one stage to another.
- `modules`: a list of modules to use in the stage.

### Modules

The modules block contains a list of modules to use in the recipe. Each module is a YAML snippet that defines a set of instructions. The common structure is:

```yml
- name: name-of-the-module
  type: type-of-the-module
  # specific fields for the module type
```

Refer to the [Use Modules](/vib/en/use-modules) article for more information on how to use modules in a recipe and [Built-in Modules](/vib/en/built-in-modules) for a list of the built-in modules and their specific fields.

You can also write your custom modules by making a Vib plugin, see the [Making a Plugin](/vib/en/making-plugin) article for more information.

### Copying files between stages

You can copy files between stages using the `copy` field. This consists of a list of files or directories to copy from another stage. Each item in the list is a YAML snippet that defines the source and destination of the copy operation. The common structure is:

```yml
- from: stage-id-to-copy-from
  paths:
    - src: /path/to/source
      dst: /path/to/destination
```

For example, to copy the `/path/to/output` directory from the `build` stage to the `/app` directory in the `dist` stage, you can use the following snippet:

```yml
- from: build
  paths:
    - src: /path/to/output
      dst: /app
```

so it becomes available in the `dist` stage.
