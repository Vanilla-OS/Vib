<div align="center">
    <img src="https://vib.vanillaos.org/assets/brand/logo/svg/full-mono-light.svg#gh-light-mode-only" height="64">
    <img src="https://vib.vanillaos.org/assets/brand/logo/svg/full-mono-dark.svg#gh-dark-mode-only" height="64">
    <p>Vib (Vanilla Image Builder) is a tool that streamlines the creation of container images. It achieves this by enabling users to define a recipe consisting of a sequence of modules, each specifying a particular action required to build the image. These actions may include installing dependencies or compiling source code.
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
vib build recipe.yml
```

this will parse the recipe.yml to a Containerfile, which can be used to build
the image with any container image builder, such as `docker` or `podman`.
