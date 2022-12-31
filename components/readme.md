# Dapr Components

The application requires the follow Dapr components for operation:

- A state store component, with a name `statestore` (this is the default name and can be changed)
- A pub/sub component, with a name `pubsub` (this is the default name and can be changed)

When working locally with Dapr these components with these names are deployed by default (i.e. when running `dapr init`), and backed by a Redis container, so no extra configuration or work is required.

When deploying to Kubernetes the Redis state provider needs to be stood up, Helm is an easy way to do this, see [deploy/readme.md](../deploy/readme.md). The Dapr Store Helm chart then will install the relevant Dapr component definitions `statestore` and `pubsub` to use this Redis instance

## Required Components

There are two Dapr components required for the app to function:

- `pubsub.yaml` - Default pub/sub component of type **pubsub.redis**, included here for reference. No configuration required.

  - _Running locally_: This should exist by default in your dapr components dir, eg. $HOME/.dapr/components
  - _Running in Kubernetes_: The daprstore Helm chart should deploy this component for you

- `statestore.yaml` - Default state component of type **state.redis**, included here for reference. No configuration required

  - _Running locally_: This should exist by default in your dapr components dir, eg. $HOME/.dapr/components
  - _Running in Kubernetes_: The daprstore Helm chart should deploy this component for you

## Optional Components

These components do not need to be deployed/installed in order for the application to run and function

- `orders-email.yaml` - Component type: **bindings.twilio.sendgrid**. Used by the orders service. Uses SendGrid to email & notify users when their order is received. You will require a SendGrid account and API key to set this up.

  - _Running locally_: Rename/copy the sample file removing .sample, and set your SendGrid API key. Then copy to your default dapr components dir, eg. `$HOME/.dapr/components`
  - _Running in Kubernetes_: Use kubectl to apply this file to your cluster and the same namespace as your application

- `orders-report.yaml` - Component type: **bindings.azure.blobstorage**. Used by the orders service. Uses Azure storage to store simple "order reports" as text files when an order is received. You will require an Azure subscription, and storage account to set this up.

  - _Running locally_: Rename/copy the sample file removing .sample, and set your storage account details & key. Then copy to your default dapr components dir, eg. `$HOME/.dapr/components`
  - _Running in Kubernetes_: Use kubectl to apply this file to your cluster and the same namespace as your application

## Component Names

The names of the components can be changed if you wish, using the following env vars, but this is not generally not recommended or required:

- `DAPR_PUBSUB_NAME` - Name of the Dapr pub/sub component to use for orders. Default is `pubsub`
- `DAPR_EMAIL_NAME` - Name of the Dapr SendGrid component to use for sending order emails. Default is `orders-email`
- `DAPR_REPORT_NAME` - Name of the Dapr Azure Blob component to use for saving order reports. Default is `orders-report`
