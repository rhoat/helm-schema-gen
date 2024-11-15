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
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name:

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths: []
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}
```

Now if you use the plugin and pass the `values.yaml` to it, you will
get the JSON Schema for the `values.yaml`

```json
$ helm schema-gen values.yaml
{
    "$schema": "http://json-schema.org/schema#",
    "type": "object",
    "properties": {
        "affinity": {
            "type": "object"
        },
        "fullnameOverride": {
            "type": "string"
        },
        "image": {
            "type": "object",
            "properties": {
                "pullPolicy": {
                    "type": "string"
                },
                "repository": {
                    "type": "string"
                }
            }
        },
        "imagePullSecrets": {
            "type": "array"
        },
        "ingress": {
            "type": "object",
            "properties": {
                "annotations": {
                    "type": "object"
                },
                "enabled": {
                    "type": "boolean"
                },
                "hosts": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "host": {
                                "type": "string"
                            },
                            "paths": {
                                "type": "array"
                            }
                        }
                    }
                },
                "tls": {
                    "type": "array"
                }
            }
        },
        "nameOverride": {
            "type": "string"
        },
        "nodeSelector": {
            "type": "object"
        },
        "podSecurityContext": {
            "type": "object"
        },
        "replicaCount": {
            "type": "integer"
        },
        "resources": {
            "type": "object"
        },
        "securityContext": {
            "type": "object"
        },
        "service": {
            "type": "object",
            "properties": {
                "port": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "serviceAccount": {
            "type": "object",
            "properties": {
                "create": {
                    "type": "boolean"
                },
                "name": {
                    "type": "null"
                }
            }
        },
        "tolerations": {
            "type": "array"
        }
    }
}
```

You can save it to a file like this

```
$ helm schema-gen values.yaml > values.schema.json
```

## Issues? Feature Requests? Proposals? Feedback?

Put them all in [GitHub issues](https://github.com/Rhoat/helm-schema-gen/issues) 😁
