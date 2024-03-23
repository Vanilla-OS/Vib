---
Title: How to use Vib modules
Description: How to use predefined and custom modules in your Vib recipes.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
Tags:
  - modules
---

Modules are a fundamental part of Vib, likely the thing you will interact with the most. We can see them as the building blocks of your container image, each one performing a specific task.

## Familiarize with Vib Recipes

Before diving into the modules, it's important to understand the structure of a Vib recipe.

![Vib Recipe Structure](https://vib.vanillaos.org/uploads/vib-recipe-structure.png)

As you can see, a recipe has two main entities:

- **Recipe**: your recipe, the YAML file that contains the instructions to build the image.
- **Stages**: one or more stages that define a set of operations to perform. Each stage can have a set of modules that will be executed in order.

Think of a recipe as a floor plan for a house, and the stages as the rooms. Each room has a set of tasks to complete, and each task is a module but all the rooms together make the house, your container image.

To get more information about the structure of a recipe and its fields, please refer to the [recipe structure](/vib/en/recipe-structure) article.

## Architecture of a Module

A module is a YAML snippet that defines a set of instructions, the common structure is:

```yaml
name: name-of-the-module
type: type-of-the-module
# specific fields for the module type
```

While the `name` and `type` fields are mandatory, the specific fields depend on the module type. For example, a `shell` module has a `commands` field that contains the shell commands to execute to complete the task.

You will find that some modules has a common `source` field, this instructs Vib to download a resource required for the module to work:

```yaml
- name: vanilla-tools
  type: shell
  source:
    type: tar
    url: https://github.com/Vanilla-OS/vanilla-tools/releases/download/continuous/vanilla-tools.tar.gz
  commands:
    - mkdir -p /usr/bin
    - cp /sources/vanilla-tools/lpkg /usr/bin/lpkg
    - cp /sources/vanilla-tools/cur-gpu /usr/bin/cur-gpu
    - chmod +x /usr/bin/lpkg
    - chmod +x /usr/bin/cur-gpu
```

In the above example we define a `shell` module that downloads a tarball from a GitHub release and then copies the binaries to `/usr/bin`. A source can be of two types:

- `tar`: a tarball archive. It can also define a `Checksum` field to verify the integrity of the downloaded file using a `sha256` hash.
- `git`: a Git repository.

In the latter case, you can specify the branch, tag or commit to checkout like this:

```yaml
name: apx-gui
type: meson
source:
  type: git
  url: https://github.com/Vanilla-OS/apx-gui
  branch: main
  commit: latest
modules:
  - name: apx-gui-deps-install
    type: apt
    source:
      packages:
        - build-essential
        - gettext
        - libadwaita-1-dev
        - meson
```

Supported fields for a git source are:

- `tag`: the tag to checkout, collides with `branch` and `commit`.
- `branch`: the branch to checkout, collides with `tag`.
- `commit`: the commit to checkout, collides with `tag` and `branch`. It can be a commit hash or `latest` to checkout the latest commit.

## Built-in Modules

Vib comes with a set of predefined modules that you can use in your recipes. You can find the list of available modules in the [list of modules](/vib/en/built-in-modules) article.

## Custom Modules via Plugins

You can also extend Vib with custom modules by writing a plugin. Please refer to [making a plugin](/vib/en/make-plugin) for more information.
