# OCF Scheduler CF CLI Plugin

# What is OCFScheduler

**ocf-scheduler-cf-plugin** enables you to interact with the OCF Scheduler API for scheduling calls and jobs.

## How to install

Clone the repo, build the binary, and install the plugin:
```
git clone https://github.com/cloudfoundry-community/ocf-scheduler-cf-plugin
cd ocf-sceduler-cf-plugin
make install
```

## Extended Build Metadata

The extended build metadata is available by executing the plugin itself.

```
./ocf-scheduler-cf-plugin
```

## Release Engineering Information

### GMake targets for pipelining

Generate release artifacts in the ***releases*** directory.  Sha256 checksum files are
created for each artifact as well.
```
gmake ci-release VERSION=<major.minor,patch> [SEMVER_PRERELEASE=<prerelease-metadata>] [SEMVER_BUILDMETA=<buid-metadata>]
```

Cleans up after the artifacts are copied 
```
gmake release-clean

```
### GMake targets for manual release engineering

You can use the pipelining targets, but for manual release engineering testing the following
additional targets can be used for release build testing.

Same as ci-release but the version information is picked up from the latest tag.  The patch
version is incremented to avoid version confusion for actual releases and the prerelease is set to dev.  
The same GMake variables used on ci-release can be used on release-all
```
gmake release-all
```

## GMake targets for development engineering

In addtion to the release engineering targets, the following target will build the artifact for your local
environment.   The built artifact is located in the repo top directory.  The install target will build and install
the plugin into your local cf command. The same GMake variables for release engineering are available for development
engineering.
```
gmake build
./ocf-scheduler-cf-plugin
```
or
```
gmake install
./ocf-scheduler-cf-plugin
```

Clean up the build or install target generated artifacts
```
gmake clean
```

## References

[CF Plugin API](https://github.com/cloudfoundry/cli/blob/master/plugin/plugin_examples/DOC.md) documentation

## History

* v0.0.3 - Now with Job disk/memory quotas
* v0.0.2 - Now with a version (but not aversion)
* v0.0.1 - Initial release
