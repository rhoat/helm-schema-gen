# helm schema gen plugin

![](https://github.com/Rhoat/helm-schema-gen/workflows/goreleaser/badge.svg)

So that you don't have to write values.schema.json by hand from scratch for your Helm 3 charts

[Helm](https://helm.sh) plugin to generate [JSON Schema for values yaml](https://helm.sh/docs/topics/charts/#schema-files)

## Note about the project

This is a fork of [karuppiah7890](https://github.com/karuppiah7890/helm-schema-gen)'s project. Their work laid the foundation for what this project is today and what it will become. We appreciate their contributions and commitment to the project.

## Install

The plugin works with both Helm v2 and v3 versions as it's agnostic to the Helm
binary version

```
$ helm plugin install https://github.com/rhoat/helm-schema-gen
Downloading https://github.com/rhoat/helm-schema-gen/releases/download/0.0.8/helm-schema-gen-linux-amd64.tgz
Preparing to install into /home/user/.local/share/helm/plugins/helm-schema-gen
Installed plugin: schema-gen
```

### Offline install
To install this offline:
- clone the repository
- copy it to the desired host
- cd to that directory
- run `helm plugin install .`

But note that the schema feature is present only in Helm v3 charts, so Helm
chart still has to be v3, meaning - based on the Helm chart v3 spec. And the
schema validation is only done in Helm v3. Read more in the
[Schema Files](https://helm.sh/docs/topics/charts/#schema-files) section of the
Helm official docs.

## Usage

The plugin works with both Helm v2 and v3 versions

Let's take a sample `values.yaml` like the below

```yaml
autoscaling:
  enabled: false
  minReplicas: 1 
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80
```

Now if you use the plugin and pass the `values.yaml` to it, you will
get the JSON Schema for the `values.yaml`

```json
{
    "$schema": "http://json-schema.org/schema#",
    "type": "object",
    "properties": {
        "autoscaling": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean"
                },
                "maxReplicas": {
                    "type": "integer"
                },
                "minReplicas": {
                    "type": "integer"
                },
                "targetCPUUtilizationPercentage": {
                    "type": "integer"
                }
            }
        }
    }
}

```

You can save it to a file like this

```bash
helm schema-gen generate ./values.yaml --destination json.schema
```

### schemagen helpers
> [!WARNING]  
> This is an experimental feature.

For more advanced use cases, if you want to modify the type that is output you can leverage the schemagen comments.
> [!TIP]
> The tool looks at the Headcomment on the yaml node, in order to make sure your comment isn't picked up as a footer its best to leave a space above your comment shown below. 
```yaml
autoscaling:
  enabled: false
  minReplicas: 1 

  # +schemagen:type:number
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80
```
When run with the `-schemagen` flag, you will be able to specify the type that is output in your json. 

> [!NOTE]
> There is no restriction on the type that can be put in here, so if you mistype something or put in garbage you will get garbage out.

```bash
helm schema-gen generate ./testdata/values.yaml --schemagen --destination json.schema
```

```
{
    "$schema": "http://json-schema.org/schema#",
    "type": "object",
    "properties": {
        "autoscaling": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean"
                },
                "maxReplicas": {
                    "type": "number"
                },
                "minReplicas": {
                    "type": "integer"
                },
                "targetCPUUtilizationPercentage": {
                    "type": "integer"
                }
            }
        }
    }
}
```

## Issues? Feature Requests? Proposals? Feedback?

Put them all in [GitHub issues](https://github.com/Rhoat/helm-schema-gen/issues) 😁
