# argocd-offline-cli

A small [Argo CD](https://argo-cd.readthedocs.io/en/stable/) CLI utility, built on top of Argo CD Go packages, that can be used to preview the Kubernetes resources generated from an [ApplicationSet](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/applicationset-specification/), without the need to connect to an actual Argo CD server.

## Requirements

* A recent version of [Helm v3](https://helm.sh/).

## Limitations

Only a few [generators](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/Generators/) are supported:
- **List** generator
- **Matrix** generator
- **Merge** generator
- **Git** generator (with local repositories only, see below)

Only Helm source repositories are supported for manifest generation.

## Usage

### Configuration

The `HELM_REPO_USERNAME` and `HELM_REPO_PASSWORD` environment variables can be specified in order to provide the default credentials that should be used to authenticate to Helm repositories. If not specified, the local `helm` command settings may be used to authenticate (if present).

### Preview Application(s) from an ApplicationSet

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest
```

#### Example: filter by application name

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest -n app-name
```

#### Example: filter by application name, display in yaml format

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest -n app-name -o yaml
```

### Using the Git Generator with Local Repositories

Since this tool is designed to work offline, the Git generator cannot fetch from remote repositories. Instead, you can map remote repository URLs to local directories using the `--local-repo` flag.

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest \
  --local-repo "https://github.com/org/repo.git=/path/to/local/repo"
```

You can specify multiple mappings:

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest \
  -l "https://github.com/org/repo1.git=/path/to/local/repo1" \
  -l "https://github.com/org/repo2.git=/path/to/local/repo2"
```

The Git generator supports both directory and file-based generation from the local repositories.

### Preview Resource manifest(s) from an ApplicationSet

```shell
argocd-offline-cli appset preview-resources /path/to/application-set-manifest
```
