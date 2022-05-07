# OCF Scheduler CF CLI Plugin

# What is cf-plugin-deploy?

**cf-plugin-ocf-scheduler** allows you to specify initial users and their credentials, their role/space associations, as well as limiting what brokers are accessible from which spaces at the time of the CF deployment.

## How to install

Clone the repo, build the binary, and install the plugin:
```
git clone https://github.com/starkandwayne/cf-plugin-ocf-scheduler
cd cf-plugin-ocf-scheduler
make
make install
```

## References

[CF Plugin API](https://github.com/cloudfoundry/cli/blob/master/plugin/plugin_examples/DOC.md) documentation

