# Abbreviated Helm Example

This base helm operator also implements what's required to support disconnected
environments. It's abbreviated, but here's the gist.

- The [deployment manifest](./helm-charts/nginx/templates/deployment.yaml#L34)
  has a variable reference `.Values.image.nginx`.
  
- This location typically refers to a value in
  [values.yaml](./helm-charts/nginx/values.yaml) but for our purposes, we don't
  want to set this value there.

- Instead, we set this value to be overridden by our
  [watches.yaml](./watches.yaml) which is reading the environment variable
  `RELATED_IMAGE_NGINX`.

- We then update our [controller manager
  deployment](./config/manager/manager.yaml#L75) spec to set our desired related
  image by digest or by tag. If you do it by tag, it gets pinned at bundle
  creation. 

- Then you change your Makefile to add image digest pinning flags to
  `operator-sdk generate bundle`. This step is not reflected in this repo, but
  it's the same as the Go example.