# Crisscross

Crisscross is a special [Crossplane](https://crossplane.io/) provider that
allows you to write "Providers as a Function" (PaaF).

## How It Works

Traditional Crossplane providers utilize
[crossplane-runtime](https://github.com/crossplane/crossplane-runtime) to defer
much of the logic for interacting with Kubernetes objects to the generic managed
resource reconciler it implements. You can find more documentation on this model
[here](https://crossplane.io/docs/v0.14/contributing/provider_development_guide.html),
or check out one of the many providers in the Crossplane community. This is a
powerful model that has allowed many community members to create new providers
in a few hours. However, though much of the code is boilerplate from a template
such as [provider-template](https://github.com/crossplane/provider-template),
there is still quite a lot of machinery required to add support for a simple
API. This can be especially challenging if you are not familiar with Kubernetes
controllers.

Crisscross shifts more of the burden that currently resides on provider authors
to a single system. Instead of provider authors implementing many managed
reconcilers for their different API types, they deploy a single service and a
corresponding `Registration` object in their cluster for each API type. When the
`Registration` is created, Crisscross spins up a new controller watching for the
referenced API type that will call out to the service at the supplied endpoint.

For example, a `Registration` for a `Bucket` on GCP could look as follows:

```yaml
apiVersion: crisscross.crossplane.io/v1alpha1
kind: Registration
metadata:
  name: bucket-paaf
spec:
  typeRef:
    apiVersion: storage.gcp.crossplane.io/v1alpha3
    kind: Bucket
  endpoint: http://172.18.0.2:32062
```

The `spec.endpoint` field can point to any service, whether it is a traditional
`Pod`, a [Knative](https://knative.dev/) `Service`, or a public API on the
internet. The only requirement for this service is that it serves the following
methods:

- `/observe`
- `/create`
- `/update`
- `/delete`

Crisscross will send the managed resource to these endpoints, and will take
appropriate action in the cluster based on the response. For an example of how
simple a "PaaF" can be, take a look at [nop-paaf](/examples/nop-paaf), which
just reports that a resource exists and all other operations are no-ops.

## License

Crisscross is under the Apache 2.0 license.
