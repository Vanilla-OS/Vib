---
Title: Project Structure
Description: How to structure your Vib project.
PublicationDate: 2024-02-13
Authors:
  - mirkobrombin
  - kbdharun
Tags:
  - project
---

Vib only requires a `vib.yml` file to build in the root of your project. However, to take full advantage of Vib, you can follow a specific project structure.

## Standard Project

A project is a directory containing a `vib.yml` file, this is the easiest way to use Vib in your existing projects, whatever their structure is. Then simply run `vib build` to build the image according to your recipe.

The following is an example of a project structure:

```plaintext
my-project/
├── vib.yml
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

Here are some details about the structure:

- `vib/` is the directory containing the Vib project.
- `includes.container/` is the directory containing the files to be included in the image. It can contain any file or directory you want to include in the image. The files in this directory will be copied to the root of the image following the same structure.
- `modules/` is the directory containing the modules used in the recipes. You can create as many modules directories as you want, naming them as you prefer. Each module directory contains one or more YAML files, each one representing a module, name them as you prefer.
- `vib.yml` is the recipe file for the image. You can have multiple `vib.yml` files in the same project, each one representing a different image. For example, you can have a `dev.yml` and a `prod.yml` file to build different images for development and production environments, then build them with `vib build dev.yml` and `vib build prod.yml`.

### Include Modules in the Recipe

You can define your modules directly in the recipe file but the above structure is recommended to keep the project organized and to reuse the modules across different recipes. So, once you have defined your modules directories, you can include them in the recipe file using the `include` module:

```yml
- name: deps-modules
  type: includes
  includes:
    - modules/node.yml
    - modules/python.yml

- name: proj-modules
  type: includes
  includes:
    - modules/myproject.yml
```

#### Remote Modules

Vib has support for remote modules, you can include them in the recipe file using the module URL or the `gh` pattern:

```yml
- name: deps-modules
  type: includes
  includes:
    - https://my-repo.com/modules/node.yml
    - gh:my-org/my-repo:branch:modules/python.yml
```

As you can see in the above example, we are explicitly including each module in the recipe file and not pointing to the whole `modules` directory. This is because the `include` module ensures each module gets included in the exact order you specify, ensuring the build process is predictable.

### Usecase of the includes.container Directory

As mentioned, the `includes.container` directory contains the files to be included in the image. This directory is useful to include files that are not part of the project, for example, configuration files, desktop files, or any other file you want to include in the image.

This is useful especially when you need to configure the Linux system with custom configuration files or new `systemd` services.

### Use the `adds` Directive

Optionally, you can use the `adds` directive to include more directories and files in the image:

```yml
adds:
  - extra-files/
  - /etc/my-config-file
```
