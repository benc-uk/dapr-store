# daprstore

![Version: 0.7.0](https://img.shields.io/badge/Version-0.7.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.7.0](https://img.shields.io/badge/AppVersion-0.7.0-informational?style=flat-square)

A reference application showcasing the use of Dapr

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| auth.clientId | string | `nil` | Set this to enable authentication, leave unset to run in demo mode |
| cart.annotations | string | `nil` | Dapr store cart annotations |
| cart.replicas | int | `1` | Dapr store cart replica count |
| daprComponents.deploy | bool | `true` | Enable to deploy the Dapr components |
| daprComponents.pubsub.name | string | `"pubsub"` | Dapr pubsub component name |
| daprComponents.pubsub.redisHost | string | `"daprstore-redis-master:6379"` | Hostname of redis, fullnameOverride should be used when deploying redis helm chart |
| daprComponents.state.name | string | `"statestore"` | Dapr state store component name |
| daprComponents.state.redisHost | string | `"daprstore-redis-master:6379"` | Hostname of redis, fullnameOverride should be used when deploying redis helm chart |
| frontendHost.annotations | string | `nil` | Dapr store frontend host annotations |
| frontendHost.replicas | int | `1` | Dapr store frontend host replica count |
| image.pullSecrets | list | `[]` | Any pullsecrets that are required to pull the image |
| image.registry | string | `"ghcr.io"` | Image registry, only change if you're using your own images |
| image.repo | string | `"benc-uk/daprstore"` | Image repository |
| image.tag | string | `"latest"` | Image tag |
| ingress.certIssuer | string | `nil` | Cert manager issuer, leave unset to run in insecure mode |
| ingress.certName | string | `nil` | Set this to enable TLS, leave unset to run in insecure mode |
| ingress.host | string | `nil` | Ingress host DNS name |
| orders.annotations | string | `nil` | Dapr store orders annotations |
| orders.replicas | int | `1` | Dapr store orders replica count |
| products.annotations | string | `nil` | Dapr store products annotations |
| products.replicas | int | `1` | Dapr store products replica count |
| resources.limits.cpu | string | `"100m"` | CPU limit for the containers, leave alone mostly |
| resources.limits.memory | string | `"200M"` | Memory limit for the containers, leave alone mostly |
| users.annotations | string | `nil` | Dapr store users annotations |
| users.replicas | int | `1` | Dapr store users replica count |

