---
Title: Project Structure
Description: How to structure your Vib project.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
Tags:
  - project
---

Vib only requires a `vib.yaml` file to build in the root of your project. However, to take full advantage of Vib, you can follow a specific project structure.

## Standard Project

A project is a directory containing a `vib.yaml` file, this is the easiest way to use Vib in your existing projects, whatever their structure is. Then simply run `vib build` to build the image according to your recipe.

The following is an example of a project structure:

```plaintext
my-project/
├── vib.yaml
```

## Vib Project

A Vib project is a directory dedicated to Vib, it can be placed in your existing project or as a standalone directory (e.g. a dedicated repository). It can contain multiple recipes to build different images. Usually, a Vib project also contains a `includes.container` directory with extra files to be included in the image and one or more directories to store the modules used in the recipes.

The following is an example of a Vib project structure:

```plaintext
my-project/
├── folder and files of your project
├── vib/
│   ├── includes.container/
│   │   ├── etc/
│   │   │   └── my-config-file
│   │   ├── usr/
│   │   │   └── share/
│   │   │       └── applications/
│   │   │           └── my-app.desktop
│   ├── modules/
│   │   ├── node.yaml
│   │   ├── python.yaml
│   │   └── myproject.yaml
│   └── vib.yaml
```

### Structure Details

Here some details about the structure:

- `vib/` is the directory containing the Vib project.
- `includes.container/` is the directory containing the files to be included in the image. It can contain any file or directory you want to include in the image. The files in this directory will be copied to the root of the image following the same structure.
- `modules/` is the directory containing the modules used in the recipes. You can create as many modules directories as you want, naming them as you prefer. Each module directory contains one or more YAML files, each one representing a module, name them as you prefer.
- `vib.yaml` is the recipe file for the image. You can have multiple `vib.yaml` files in the same project, each one representing a different image. For example, you can have a `dev.yaml` and a `prod.yaml` file to build different images for development and production environments, then build them with `vib build dev.yaml` and `vib build prod.yaml`.

### Include Modules in the Recipe

You can define your modules directly in the recipe file but the above structure is recommended to keep the project organized and to reuse the modules across different recipes. So, once you have defined your modules directories, you can include them in the recipe file using the `include` module:

```yaml
- name: deps-modules
  type: includes
  includes:
    - modules/node.yaml
    - modules/python.yaml

- name: proj-modules
  type: includes
  includes:
    - modules/myproject.yaml
```

#### Remote Modules

Vib has support for remote modules, you can include them in the recipe file using the module URL or the `gh` pattern:

```yaml
- name: deps-modules
  type: includes
  includes:
    - https://my-repo.com/modules/node.yaml
    - gh:my-org/my-repo/modules/python.yaml
```

As you can see in the above example, we are explicitly including each module in the recipe file and not pointing to the whole modules directory. This is because the `include` module ensure each module gets included in the exact order you specify, ensuring the build process is predictable.

### Usecase of the includes.container Directory

As mentioned, the `includes.container` directory contains the files to be included in the image. This directory is useful to include files that are not part of the project, for example, configuration files, desktop files, or any other file you want to include in the image.

This is useful expecially when you need to configure the Linux system with custom configuration files or new systemd services.

### Use the `adds` Directive

Optionally, you can use the `adds` directive to include more directories and files in the image:

```yaml
adds:
  - extra-files/
  - /etc/my-config-file
```
