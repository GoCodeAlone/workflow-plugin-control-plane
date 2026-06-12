# workflow-plugin-control-plane

Product-neutral Workflow control-plane descriptor, envelope, and registry
contracts.

This plugin package is a public contract lane. It validates descriptor data that
Workflow hosts can choose to render or store, but it does not authorize actions,
bind HTTP routes, dispatch providers, issue credentials, persist host state,
approve deployments, or run rollbacks.

## Packages

- `descriptors`: route/action, provider handoff, and admin contribution
  descriptor validators.
- `envelopes`: state/event/audit envelope validators with opaque handle and
  retention metadata checks.
- `registry`: descriptor provenance, digest, allowlist epoch, downgrade floor,
  revocation freshness, parser version, and trust-root generation validators.

The wire contract is v1alpha1 and uses protocol version
`control-plane.v1alpha1`.
