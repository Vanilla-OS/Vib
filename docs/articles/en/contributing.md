---
Title: Contributing
Description: We welcome contributions from the community. Learn how to contribute to the Vib project.
PublicationDate: 2024-02-14
Listed: true
Authors:
  - mirkobrombin
Tags:
  - contributing
  - development
---

We welcome contributions from the community, Vib is an Open Source project and we are always looking for new features, bug fixes, and improvements. This guide will help you get started with contributing.

## How to Contribute

There are many ways to contribute, from writing documentation and tutorials, to testing, submitting bug reports, or writing code to fix bugs or add features.

### Writing Documentation and Tutorials

If you want to contribute to the documentation, you can do so by submitting a pull request to the [Vib](https://github.com/Vanilla-OS/Vib) repository. The documentation is written in Markdown and is located in the `docs/articles` folder. We use [Chronos](https://github.com/Vanilla-OS/Chronos) to manage the documentation, so make sure to follow the article metadata structure [described here](https://github.com/Vanilla-OS/Chronos/tree/main?tab=readme-ov-file#article-structure).

Documentation should be clear, concise, and easy to understand. If you are writing a tutorial, make sure to include all the necessary steps to reproduce the tutorial, taking care of documenting terms and concepts that might not be familiar to all readers.

Provide examples and code snippets to help the reader understand the concepts you are explaining, use illustrations only to illustrate complex structures or concepts, for stuff like the structure of a folder use a code based representation instead.

### Testing and Submitting Bug Reports

If you find a bug in Vib, please submit an issue to the [Vib](https://github.com/Vanilla-OS/Vib/issues) repository. Before submitting a bug report, make sure to check if it has already been reported, and if not, provide as much information as possible to help us reproduce the issue.

Bug reports are very important to us, and we appreciate the time you take to submit them. Just make sure to report bugs in the context of the Vib project, if you are using a recipe and you find a bug in it, please report it to the recipe repository.

### Writing Code

If you are a developer and want to contribute to the Vib project by writing code, you can do so by submitting a pull request to the [GitHub repository](https://github.com/Vanilla-OS/Vib). Before writing code, make sure to check if the feature you want to implement is already being worked on, and if not, open an issue to discuss it with the maintainers or join our [Discord](https://vanillaos.org/community) to discuss it with the community, look for the `#vib` channel.

We appreciate your time and effort in contributing to the Vib project, we would be sorry if you write code that we cannot merge, so make sure to discuss your ideas before starting.

#### Extending Built-in Modules

If you want to add a new built-in module to Vib, please consider the following:

- **Is it really necessary?** Make sure that the module you want to add is not already doable with the existing modules, and that it is a common use case, if the process is too complex using the available modules, then it might be worth it. Modules should be generic and reusable, if you need a module for a specific use case, consider writing a plugin instead.
- **Does it require a new dependency?** If the module you want to add requires a new dependency, make sure that it is a widely used library, and that it is not too heavy. We want to keep Vib as lightweight as possible.
- **Use self-explanatory names**: The name of the module should be self-explanatory, and it should be as short as possible, while still being descriptive. The same applies to the module's parameters, for example, if you are writing a module to copy files (that should not be the case since a shell module is good enough for that), the parameters should be `from` and `to`, not `originalPath` and `destinationPath`, while both are correct, the first is more concise and easier to understand.
- **Distro-agnostic first**: If the module you want to add is distro-agnostic, it is more likely to be accepted, if it is not, make sure to provide a good reason for it.

#### Code of Conduct

When contributing to this project, please make sure to read and follow our [Code of Conduct](https://vanillaos.org/code-of-conduct). This document helps us to create a welcoming and inclusive community, and we expect all contributors to follow it.

#### Contribution Guidelines

Before contributing, please make sure to read our [Contribution Guidelines](https://github.com/Vanilla-OS/.github/blob/main/CONTRIBUTING.md), to learn how to write commits, submit pull requests, and more.
