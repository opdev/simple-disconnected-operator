# Abbreviated Ansible Example

This is an abbreviated example showing how to functionally enable disconnected
environments.

- The `watches.yaml` file does support an `overrideValues` key like the Helm
  operator project layout does. Ansible operators need to rely on the variables
  defined in their roles to accomplished RELATED_IMAGES.

- In our role's [defaults/main.yml](./roles/memcached/defaults/main.yml#L4-L5),
  a default value for the memcached image is stored, along with a variable that
  reads from the environment using the `lookup` plugin`. Any combination of
  these values can probably be moved to the vars/main.yml file, but for
  simplicity, I've kept them together here.

- Then, in our [task](./roles/memcached/tasks/main.yml#L28), we use Jinja
  templating for the image value of our deployment, indicating that we should
  always use the value looked up in the RELATED_IMAGE_MEMCACHED environment
  variable, and instead use our fallback/default value if that's empty.

- Finally, as with the other implementations, we update our [controller
  manager](./config/manager/manager.yaml#L78) to pass through the environment
  variable to the controller at runtime. Here, I've used a tag, but
  operator-sdk, with the appropriate flags, will pin this in our bundle. See the
  Golang example README.md for exactly how to enable that feature.