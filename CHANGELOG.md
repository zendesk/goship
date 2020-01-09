## 1.0.5 (January 9, 2020)

NEW FEATURES:
* [[GH-13](https://github.com/zendesk/goship/pull/13)]: Implemented --ssh-command flag for ssh command for remote execution


## 1.0.4 (March 15, 2019)

* EC2 provider: fallback to Private IP/DNS when public is not available


## 1.0.3 (March 14, 2019)

BUG FIXES:
* [[GH-11](https://github.com/zendesk/goship/pull/11)]: Add missing ssh/scp extra args support


## 1.0.2 (December 17, 2018)

BUG FIXES:
* [[GH-6](https://github.com/zendesk/goship/pull/6)]: AwsEc2Cache Refresh() fix
* [[GH-8](https://github.com/zendesk/goship/pull/8)]: Print errors when checking for the newer version
* [[GH-9](https://github.com/zendesk/goship/pull/9)]: Fix cache refresh messages

OTHER:
* [[GH-7](https://github.com/zendesk/goship/pull/7)]: Enable unit-tests and test PR building with Travis


## 1.0.1 (December 13, 2018)

NEW FEATURES:
* Check for the newest version when using `--verbose` flag

BUG FIXES:
* [[GH-1](https://github.com/zendesk/goship/pull/1)]: execute CheckForNewVersion after parsing flags
* Add warning message when no providers configured
* Fix joining paths when creating cache files
* Fix config file example


## 1.0.0 (December 11, 2018)

Initial public release
