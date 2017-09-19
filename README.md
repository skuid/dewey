# 🤓 dewey

Dewey is an application for syncing a list of Docker repositories for registries that don't support `v2/_catalog`. Based on a list of "orgs", Dewey will generate the same output as `v2/_catalog` for a subset of the content of the registry. This way, you only include what you care about. Also, you can include a static list of repositories that you wish add in the event that you only want 1 or 2 of the repositories for a given org.

## Supported Registries

* Dockerhub
* Quay

## Configuration

The default path is at `/opt/dewey/config.yaml`, then `$HOME/.dewey`, then `.`. The format is as follows:

```yaml
registries:
  - name: quay
    kind: quay
    password: somePassword
    orgs:
      - skuid
      - ethanfrogers
  - name: dockerhub
    kind: dockerhub
    username: someUsername
    password: somePassword
    repositories:
      - gliderlabs/consul-server            
    orgs:
      - selenium
      - jenkins
```

## Configuration oddities

Since these registries do not support `v2/_catalog`, we can't utilize Docker Registrys APIs to get our final result. This means that we need to go through application APIs to get our data.

### Credentials (Quay)

Credentials for Quay aren't the same as your `docker login` credentials. For Quay, you will need to generate a Token for an "Application".

### API Address (Dockerhub)

The API address for your service can differ from registry to registry. Dockerhub, for example, lists it's application API at `https://hub.docker.com`. Quay's is listed at `https://quay.io`. If you are using a registry hosted by some other provider, the `address` option is overridable.

## Options

```
$ ./dewey --help
Usage of ./dewey:
      --dir string        catalog file output directory (default "/opt/dewey/catalogs")
      --interval string   sync interval (default "30s")
      --pretty            pretty print output
```
