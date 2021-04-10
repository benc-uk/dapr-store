# Dapr Components

- `pubsub.yaml` - Default pub/sub component, included here for reference. No configuration required

- `statestore.yaml` - Default state component, included here for reference. No configuration required

- `orders-email.yaml` - **Optional** component for orders service. Uses SendGrid to email users when their order is received. Rename/copy the sample file removing .sample, and set your SendGrid API key. Then copy to your default dapr components dir, eg. `$HOME/.dapr/components`  
  If running in Kubernetes use kubectl to apply this file to your cluster

- `orders-report.yaml` - **Optional** component for orders service. Uses Azure storage to store data when an order is received. Rename/copy the sample file removing .sample, and set your storage account details & key. Then copy to your default dapr components dir, eg. `$HOME/.dapr/components`  
  If running in Kubernetes use kubectl to apply this file to your cluster
