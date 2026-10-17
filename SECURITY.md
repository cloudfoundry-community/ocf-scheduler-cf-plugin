# Security Policy

## Reporting a vulnerability

Please do not report security problems in public issues or pull requests.

Report them privately, in either of these ways:

- GitHub private vulnerability reporting:
  <https://github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/security/advisories/new>
- Email: <nabramovitz@fivetwenty.io>

Please include the version or commit, the steps to reproduce the problem,
and the impact you see. We will confirm that we received the report, keep
you informed while we work on a fix, and credit you in the release notes
unless you would rather not be named.

## Supported versions

Security fixes are made against the latest code on `develop`, which
becomes the next release. Please check that the problem is present
there. Earlier releases are not patched; upgrade to the release that
carries the fix.

## Scope

This policy covers the cf CLI plugin in this repository. Report problems
in the scheduler service to
[ocf-scheduler](https://github.com/cloudfoundry-community/ocf-scheduler/blob/develop/SECURITY.md).
