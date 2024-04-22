---
Title: Build using GitHub Actions
Description: How to build a Vib image using GitHub Actions.
PublicationDate: 2024-02-14
Authors:
  - mirkobrombin
  - kbdharun
Tags:
  - github
  - build
---

Many projects use GitHub to host their code, and GitHub Actions to automate their workflows. Vib can be integrated into your GitHub Actions workflow to build your images automatically. To streamline the process, you can use the [Vib GitHub Action](https://github.com/Vanilla-OS/vib-gh-action).

## Setup the Workflow

To use the Vib GitHub Action, you need to create a workflow file in the repository where your Vib recipe is located. Create a new file in the `.github/workflows` directory, for example, `vib-build.yml`, and add the following content:

```yaml
name: Vib Build

on:
  push:
    branches: ["main"]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: vanilla-os/vib-gh-action@v0.7.0
        with:
          recipe: "vib.yml"
          plugins: "org/repo:tag, org/repo:tag"

      - name: Build the Docker image
        run: docker image build -f Containerfile --tag ghcr.io/your_org/your_image:main .
```

Let's break down the workflow file:

- `name`: The name of the workflow.
- `on`: The events that trigger the workflow. In this case, the workflow runs on every push to the `main` branch and when manually triggered.
- `jobs`: A workflow can contain one or more jobs. In this case, there is only one job called `build`.
- `runs-on`: The type of machine to run the job on. In this case, the job runs on the latest version of Ubuntu; check [here](https://github.com/actions/runner-images?tab=readme-ov-file#available-images) for the available machine types.
- `steps`: The sequence of tasks to run in the job.
  - `actions/checkout@v4`: A standard action to check out the repository.
  - `vanilla-os/vib-gh-action@v0.7.0`: The Vib GitHub Action to build the image. The `with` section specifies the recipe file and additional plugins to use.
  - `run`: Contains a standard command to build the Docker image. The `--tag` option specifies the name and tag of the image, in this case, the tag is `ghcr.io/your_org/your_image:main`, you can change it according to your needs.

### Using Custom Plugins

If you are using custom Vib plugins in your recipe, you can include them in the workflow file. For example, if your plugin is named `my-plugin`, you can add the following step to the workflow file:

```yaml
# other steps
- uses: vanilla-os/vib-gh-action@v0.7.0
  with:
    recipe: "vib.yml"
    plugins: "your_org/my-plugin:v0.0.1"
# the rest of the workflow
```

The syntax `your_org/my-plugin:v0.0.1` means:

- `your_org`: The GitHub organization or user that owns the plugin.
- `my-plugin`: The name of the plugin which is the same as the repository name.
- `v0.0.1`: The version of the plugin to use, which must be a valid tag.

To use more than one plugin, simply separate them with a comma:

```yaml
# other steps
- uses: vanilla-os/vib-gh-action@v0.7.0
  with:
    recipe: "vib.yml"
    plugins: "your_org/my-plugin:v0.0.1, another_org/another-plugin:v1.2.3"
# the rest of the workflow
```

## Publish the Image to GitHub Container Registry (GHCR)

The workflow file builds the Docker image to ensure everything is working as expected. If you want to publish the image to the GitHub Container Registry (GHCR), you can rework the workflow file as follows:

```yaml
name: Vib Build

on:
  push:
    branches: ["main"]
  workflow_dispatch:

env:
  REGISTRY_USER: ${{ github.actor }}
  REGISTRY_PASSWORD: ${{ secrets.GITHUB_TOKEN }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: vanilla-os/vib-gh-action@v0.7.0
        with:
          recipe: "vib.yml"

      - name: Build the Docker image
        run: docker image build -f Containerfile --tag ghcr.io/your_org/your_image:main .

      # Push the image to GHCR (Image Registry)
      - name: Push To GHCR
        if: github.repository == 'your_org/your_repo'
        run: |
          docker login ghcr.io -u ${{ env.REGISTRY_USER }} -p ${{ env.REGISTRY_PASSWORD }}
          docker image push "ghcr.io/your_org/your_image:main"
```

In this case, the `REGISTRY_USER` and `REGISTRY_PASSWORD` environment variables are set to the GitHub actor and the GitHub token, respectively. The `docker login` command uses these credentials to authenticate with GHCR, and the `docker image push` command pushes the image to the registry.
