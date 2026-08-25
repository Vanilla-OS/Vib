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
