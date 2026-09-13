<div align="center">
    <img src="logo/svg/full-mono-dark.svg#gh-light-mode-only" height="64">
    <img src="logo/svg/full-mono-light.svg#gh-dark-mode-only" height="64">
    <p>Vib (Vanilla Image Builder) creates container images from YAML recipes. A recipe contains modules that install packages, copy files, build source code, and perform other image construction steps.
</p>
    <hr />
</div>

## Links

- [Website](https://vib.vanillaos.org/)
- [Documentation](https://docs.vanillaos.org/collections/vib)
- [Examples](https://vib.vanillaos.org/examples)

## Usage

To build an image using a recipe, you can use the `vib` command:

```sh
vib build --output Containerfile recipe.yml
```

This writes a Containerfile that can be built with Docker, Podman, or another
compatible image builder. Use `vib test recipe.yml` to validate a recipe and
`vib compile --runtime docker recipe.yml` to generate and build the image in
one command.

## Use of Generative AI

Maintainers may use generative AI tools as assistants while working on Vib. Non-trivial assisted commits disclose the tool, model, and scope of the work.

AI tools may assist with code comments, documentation, repetitive code, and issue triage. Maintainers make project decisions and review every assisted change before it is merged.

Use these trailers for non-trivial assisted commits:

```plain
Assisted-by: <tool>:<model-version>
AI-Scope: <what the tool generated and the prompt or a short prompt summary>
```

Single-line completions, renames, and formatting changes do not need trailers.

Coding agents must also follow [AGENTS.md](AGENTS.md) before changing files,
creating commits, or opening pull requests.
