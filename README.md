# Simple Disconnected Operator

This repository contains a simple operator that is intended to be friendly
towards deployment in environments with restricted networks. It's a concrete
example of the suggestions listed in [OpenShift's
documentation](https://docs.openshift.com/container-platform/4.11/operators/operator_sdk/osdk-generating-csvs.html#olm-enabling-operator-for-restricted-network_osdk-generating-csvs)
for restricted-network-friendly Golang operators.

This operator manages a custom resource called **DisconnectedFriendlyApp**. When
this resource is created, it spins up two application deployments ("busybox",
and "sleeper"). Their functionality is irrelevant, but what is important is that
the controller rolls out these deployments with image values influenced by
environment variables as recommended in the aformentioned documentation.

## I want to try this out

Because this is a demo repo, there's no catalog that contains it. That said, you
can run this using a prebuilt bundle image:

```shell
operator-sdk run bundle quay.io/opdev/simple-disconnected-operator:v0.0.1
```

Create the sample resource and then inspect the custom resource instance:

```shell
oc apply config/samples/tools_v1alpha1_disconnectedfriendlyapp.yaml
oc get dfa
```

When you're done, clean it out of your cluster

```shell
operator-sdk cleanup simple-disconnected-operator
```

## The Implementation Goals

### At a high level

Aside from avoiding access to the internet, one of the key requirements for
operators that need to run in restricted networks is that they must make
references to any image-based workloads by its digest.

Restricted Network OpenShift installations will generally contain a mirrored
repository of images that **are** accessible within the cluster. Image tags can
change the image to which they reference at any time, but a network-restricted
cluster cannot see those changes in real time.

This drives the need to ensure that any workloads that require images (e.g.
deployments, statefulsets, pods) actually get the intended image, and that the
tag hasn't changed unexpectedly between the time the cluster administrators
mirrored images into the cluster, and the time the user applied the workload.

### For Operators

Operator projects that are installed via OLM may need to create any number of
workloads in a cluster. In order to make sure all images necessary for that
operator to function properly at the time the network "plug" is disconnected,
the registry mirroring processes must know what to make available locally.

This is the motivation behind the ClusterServiceVersion's "Related Images"
metadata. This information that the operator developer(s) can provide to OLM,
letting it know exactly what images the operator intends to use, and to pull
those down and make them available.

NOTE: There are other aspects of this relationship between related images and
the mirrored registry, such as ImageContentSourcePolicy, that I'm not diving
into here as it's not directly related to this repo's goals.

### So they're available, now what?

In your operator projects, you may define your resource manifests using objects,
structs, etc. in your language. In these structs, you might define an image
value that is a hard-coded string. For the most part, this is fine... so long as
the string itself matches exactly what's in your ClusterServiceVersion's Related
Images metadata.

Having it hard-coded in your source, however, makes things tricky to keep in
sync, as you have multiple places (ClusterServiceVersion, and your source code)
where you have to update these values.

### The recommended practice

... is simply to use an environment variable and pass that into your controller,
containing a known key (to your controller) with a value that is the image
reference to use for that particular secondary resource. Your controller logic
can then pass this into your resource structs (e.g. your deployment's image spec
for your secondary resources).

### Why this is recommended

This makes it so you only have a single place to update (your
ClusterServiceVersion's manager pod enviromment variables), and your code
doesn't have to change when you want to use a different image for a given
secondary resource.

In addition, OperatorSDK has a flag to the `bundle generate` command that
automatically pins your image digests when generating your bundle. There are a
few ways that this pinning works, but the easy way (and what's demonstrated in
this repo) is to have your environment variables that need to be pinned start
with `RELATED_IMAGE_`.

So as an example `RELATED_IMAGE_MARIADB` might contain a tag-referenced to a
MariaDB container image. OperatorSDK's `bundle generate` logic will, given the
correct flags are called, find this and pin it to its digest **at the time the
bundle is generated**.

In addition, this same digest pinning logic reviews all of your
ClusterServiceVersion's deployment images (e.g. for controllers, their sidecars,
proxies, initContainers, etc.) and pins those as well.

All pinned images are then added to the Related Images section of your
ClusterServiceVersion for you.

## The concrete example

### Detecting Images Defined in the Environment

I've implemented an [Image package](./image/image.go) to handle the resolution
of the images each deployment will use. You can implement this however you want,
but in this example, I [hard-code a value to use as a
fallback](./image/image.go#L6), both prefer the environment variable if defined.

### Consuming the images in controllers

I spin up two deployments for each DisconnectedFriendlyApp resource instance -
one called busybox and one called sleeper. Reconciliation logic is fairly basic:
Create the deployment if it doesn't exist, or update it if it does and make sure
we're in sync.

The key component here is that our image value for these deployments calls our
image library to resolve what our image should be. [Here's the busybox
example](./controllers/busybox_deployment_reconciler.go#L80).

### Adding environment variables to your controller manager's runtime environment

Your controller manager is what, technically, needs to know what images it
should be using. At least, that's the case in this example.

The ClusterServiceVersion that's shipped in your operator bundle contains a
direct reference to the deployment spec that's used to run your controller
manager. This is how OLM knows what to run.

Your controller manager's deployment is merged into the CSV at bundle build time
using the kustomize tooling that's generated for you by operator SDK. So, to
that end, all you need to do is update your `config/manager/manager.yaml` to
include references to the environment variables you want to use in your code.

For this repo, I wanted to decouple the environment variable values my operator
expected vs. the environment variables that the digest pinning logic expects
(e.g. prefixed with `RELATED_IMAGE_`). [My
manager](./config/manager/manager.yaml#L76) has a `RELATED_IMAGE_FEDORA`
environment variable defined, which I then pass to the [two environment variables](./image/image.go#L8)
my controller logic actually uses (`DFA_SLEEPER_IMAGE`, `DFA_BUSYBOX_IMAGE`).

You'll notice that my `RELATED_IMAGE_FEDORA` value uses a tag reference! Here it
uses `latest` which doesn't make for the best example, but at the time of this
writing `latest` points to `quay.io/fedora/fedora:37`.

It's okay to use a tag here. We will configure bundle generation logic to
resolve this in the bundle it generates.

### Generating the bundle with pinned values

The [OpenShift
documentation](https://docs.openshift.com/container-platform/4.11/operators/operator_sdk/osdk-generating-csvs.html#olm-enabling-operator-for-restricted-network_osdk-generating-csvs)
goes through how to adjust your bundle generation flags to properly pin your
digests in the 3rd item (unfortunately I can't seem to link directly to it).

I won't go into it too much, but suffice to say that what you're adding is the
`--use-image-digests` flag, which will scan your manifest files and try to
replace any relevant image tags with digests.

Check out the [Makefile](./Makefile#L44-L47) to see how I modified mine to make
this work, as I found small issues with the documented process in the link
above, potentially due to my [OperatorSDK version](./OPERATOR_SDK_VERSION).

Once this runs (specifically, the `make bundle` target), I noticed that:

- my [controller manager
  image](./bundle/manifests/simple-disconnected-operator.clusterserviceversion.yaml#L180)
  was pinned (despite tags being used in the manager manifests. Note that you
  will need to have pushed your controller by this point for this resolution to
  work.)

- my [kube-rbac-proxy
  image](./bundle/manifests/simple-disconnected-operator.clusterserviceversion.yaml#L149)
  was pinned (despite the tag being hard-coded in what was autogenerated for
  you)

- my `RELATED_IMAGE_FEDORA` [environment variable was
  pinned](./bundle/manifests/simple-disconnected-operator.clusterserviceversion.yaml#L175)

- my CSV's relatedImages section contained all of these [image
  references](./bundle/manifests/simple-disconnected-operator.clusterserviceversion.yaml#L268)

### Confirm that all of this works

If you deploy the operator using the bundle generated from all of this work, and
then deploy the sample DisconnectedFriendlyApp instance, you'll notice that the
image swap based on environment works as expected by inspecting the status of
the instance, or just listing out the instances

```shell
$ oc get dfa
NAME                             RESOLVED BUSYBOX IMAGE                                                                          DEFAULT BUSYBOX IMAGE      RESOLVED SLEEPER IMAGE                                                                          DEFAULT SLEEPERIMAGE
disconnectedfriendlyapp-sample   quay.io/fedora/fedora@sha256:0f65bee641e821f8118acafb44c2f8fe30c2fc6b9a2b3729c0660376391aa117   docker.io/busybox:latest   quay.io/fedora/fedora@sha256:0f65bee641e821f8118acafb44c2f8fe30c2fc6b9a2b3729c0660376391aa117   docker.io/ubuntu:latest
```

### Making sure you add the infrastructure-feature annotation

The final piece is to make sure you add the
[infrastructure-feature](./config/manifests/bases/simple-disconnected-operator.clusterserviceversion.yaml#L7)
annotation to your CSV.
