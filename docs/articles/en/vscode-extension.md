---
Title: Working with Vib in Visual Studio Code
Description: Learn how to work with Vib recipes in Visual Studio Code using our extension.
PublicationDate: 2024-02-14
Authors:
  - mirkobrombin
  - kbdharun
Tags:
  - development
  - vscode
---

Visual Studio Code is a popular code editor that provides a wide range of features to help you write, debug, and deploy your code, other than being highly customizable, it also offers a wide range of extensions to enhance your development experience, for example for working with YAML files.

Vib recipes are written in YAML, and usually, a standard text editor or the YAML support provided by Visual Studio Code is enough to work with them. However, we have developed a dedicated extension for Visual Studio Code to make working with Vib recipes even easier and more efficient.

## Features

> **Note**: The Vib extension is in its early stages, and we are working to add more features and improvements. If you have any feedback or suggestions, please let us know by opening an issue on the [vib-vscode-ext](https://github.com/Vanilla-OS/vib-vscode-ext) repository.

The following features are currently available in the Vib extension (version 1.1.0):

- **Metadata validation**: checks the metadata of the recipe.
- **Modules import**: checks if the paths of a `includes` module are correct.
- **Modules name collision**: checks if the names of the modules are unique.
- **Modules type auto-completion**: suggests the type of the module to use, it works with both built-in and custom modules, for the latter it refers to the content of the `plugins` folder.

## Installation

To install the Vib extension, follow these steps:

1. Open Visual Studio Code.
2. Go to the Extensions view by clicking on the Extensions icon on the bar on the side of the window or by pressing `Ctrl+Shift+X`.
3. Search for `Vib` in the Extensions view search box.
4. Click on the `Install` button for the Vib extension.

## Usage

Once the extension gets installed, you can start using it to work with Vib recipes. For it to work, you need to put the following header at the beginning of your recipe:

```yml
# vib
```

This header is used to identify the file as a Vib recipe and to enable the extension features. In the future, we plan to have support for our dedicated file extension, but for now, this is the way to go.
