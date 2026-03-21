# Changelog

## [1.1.0](https://github.com/yavasura/terraform-provider-pbs/compare/v1.0.0...v1.1.0) (2026-03-21)


### Features

* Add complete job management, notification system, and test reorganization ([ed68739](https://github.com/yavasura/terraform-provider-pbs/commit/ed687395025858c1c16483c33f6698e6a34f3e50))
* Add PBS ISO-based Vagrant box approach ([b431cd5](https://github.com/yavasura/terraform-provider-pbs/commit/b431cd5717227a90ebdbefdb343949412b9129eb))
* Complete Go to HCL test conversion with skip markers and CI updates ([0c2c02b](https://github.com/yavasura/terraform-provider-pbs/commit/0c2c02bb225934d49b11a37eebeec09f27884e4d))
* implement comprehensive data source support for PBS resources (Phases 1-4) ([8f23dba](https://github.com/yavasura/terraform-provider-pbs/commit/8f23dba8b8d500d499158bb093fe1f6a87b7eb21))
* Refactor immutability test and add S3 provider HCL tests ([c48bf4d](https://github.com/yavasura/terraform-provider-pbs/commit/c48bf4daa53549d348c4eaaf512506b1e8fb0710))


### Bug Fixes

* add mutex to prevent datastore lock contention ([265c0b8](https://github.com/yavasura/terraform-provider-pbs/commit/265c0b8014c53c659b911894166e234acf2519ac))
* Add resource index [0] for count-based job resources ([9138cf5](https://github.com/yavasura/terraform-provider-pbs/commit/9138cf5e2b538595e77384e77fb770b4eb885a1e))
* add retry logic for reading datastore after creation ([0abeef3](https://github.com/yavasura/terraform-provider-pbs/commit/0abeef3519ecbbde80e613dbdc6e478d5d04c0fd))
* add terraform init before terraform test in CI and helper script ([2fddecf](https://github.com/yavasura/terraform-provider-pbs/commit/2fddecf5a0c76ad03ec35d37a33dd97c4d8ee346))
* Clean up formatting and remove duplicate test files ([5d7946b](https://github.com/yavasura/terraform-provider-pbs/commit/5d7946b6a03ab3b314bcc51eccafe2146d4577a2))
* Correct datasources_s3_endpoints test ([8a02f0f](https://github.com/yavasura/terraform-provider-pbs/commit/8a02f0f06e109165c4eba648f288268cad8954b8))
* Correct S3 endpoint datasource test to use id instead of name ([6473e98](https://github.com/yavasura/terraform-provider-pbs/commit/6473e98c3e6ae3d89106673ce1a491f381c5e582))
* correct terraform test command syntax - remove unsupported -chdir flag ([f6e5906](https://github.com/yavasura/terraform-provider-pbs/commit/f6e5906d8049373a217d1afaec36c7b70a0a484d)), closes [#17](https://github.com/yavasura/terraform-provider-pbs/issues/17)
* **datastore:** Make backend fields immutable and exclude from updates (fixes [#18](https://github.com/yavasura/terraform-provider-pbs/issues/18)) ([b5747ed](https://github.com/yavasura/terraform-provider-pbs/commit/b5747ed6bcf382af5698dd25a11d2193d8770967))
* Force DKMS to build and verify ZFS modules for running kernel ([9e7dc1b](https://github.com/yavasura/terraform-provider-pbs/commit/9e7dc1b3ff375799cd23c82bf662b64e1bea6f51))
* Generate test ID with timestamp in shell script instead of YAML ([31e7b8e](https://github.com/yavasura/terraform-provider-pbs/commit/31e7b8e2d74ab43dcf7a6bdd2505715856adc3cd))
* Handle default enable=true value in metrics_server datasource ([7bae03d](https://github.com/yavasura/terraform-provider-pbs/commit/7bae03db5364bb4dc3223a3ae7d287e81e41ea3e))
* Handle field deletion in notification matcher updates ([bf1160d](https://github.com/yavasura/terraform-provider-pbs/commit/bf1160d5183dc2b9aba7393baf5b72d735c5dd02))
* Handle kernel version mismatch for ZFS module loading ([9541c67](https://github.com/yavasura/terraform-provider-pbs/commit/9541c67b3dcf24545997196b62a73c6d7126c4f2))
* Handle VM reboot without conflicting with Vagrant ([4891d04](https://github.com/yavasura/terraform-provider-pbs/commit/4891d04537a9a8a12df6a225e2f0b1029610f67e))
* Implement two-stage provisioning with kernel upgrade and reboot ([99264ca](https://github.com/yavasura/terraform-provider-pbs/commit/99264ca508b7e0f44204cefedda2b994e7d065d3))
* Implement two-stage provisioning with kernel upgrade and reboot ([edc41ed](https://github.com/yavasura/terraform-provider-pbs/commit/edc41ed295abb072da156db5a96da546837f1c78))
* Improve S3 endpoint deletion and test dependency handling ([ddc1982](https://github.com/yavasura/terraform-provider-pbs/commit/ddc1982e37aa7c6e133ace82ac168653ebaec961))
* increase datastore creation retry limits for slower CI environments ([1e14a42](https://github.com/yavasura/terraform-provider-pbs/commit/1e14a423e3ffcbb54c0b3815dac2928f3e7164b2))
* Install ZFS before system upgrade to ensure module compatibility ([13299d8](https://github.com/yavasura/terraform-provider-pbs/commit/13299d8094d01df9f0ac8ac92185d5bfffe04782))
* Move metrics tests to separate top-level directories ([4b879aa](https://github.com/yavasura/terraform-provider-pbs/commit/4b879aa2fb593c94ea8b6e4f7db56869b83229f8))
* Move notification tests to separate top-level directories ([9df019c](https://github.com/yavasura/terraform-provider-pbs/commit/9df019c7fd0d4fa18161fc0c1026e7987b54eaa8))
* Move S3 immutability test to separate directory ([e37c5b2](https://github.com/yavasura/terraform-provider-pbs/commit/e37c5b23d49ae56312b8e1e1d17cbefc83d4a8dc))
* refactor mutex usage and improve retry logic for datastore operations ([fa45d9d](https://github.com/yavasura/terraform-provider-pbs/commit/fa45d9dc64b6c166397e4991c22d179ed76e1b12))
* Relax remote digest assertion to allow empty string ([d3824e2](https://github.com/yavasura/terraform-provider-pbs/commit/d3824e262547c9c57cbeb880a4bacc52b1491ef8))
* Remove duplicate test file and add timestamp to test IDs ([5f85031](https://github.com/yavasura/terraform-provider-pbs/commit/5f85031b50b6fb97ed3d64c07b2687a3cf51c502))
* Remove hardcoded PBS endpoint from immutability test ([7ec293a](https://github.com/yavasura/terraform-provider-pbs/commit/7ec293ace8ee9dfea4b9bcc1a4b13ecc7fe30461))
* Remove https:// protocol from S3 endpoint URL ([447dab9](https://github.com/yavasura/terraform-provider-pbs/commit/447dab91b82fb329af5ff5f2a6aa1790d302db2e))
* Remove invalid variable assignments from datasource tests ([ad1c461](https://github.com/yavasura/terraform-provider-pbs/commit/ad1c4613a4dbbb7a7479c1d7b8521ed91fe11128))
* Remove invalid zero values for prune job retention parameters ([2816724](https://github.com/yavasura/terraform-provider-pbs/commit/2816724934e09f63451b0edd15f63a120b3e5bc8))
* Reorganize jobs tests into subdirectories to avoid duplicate config ([2a8cdf3](https://github.com/yavasura/terraform-provider-pbs/commit/2a8cdf3732fca3eba9874b3744237924337f4265))
* replace Go log package with tflog for proper debug logging ([68e3343](https://github.com/yavasura/terraform-provider-pbs/commit/68e3343917bebfdbebcdeca5076af0419bff711c))
* Resolve lint error and test path issues ([9b8729d](https://github.com/yavasura/terraform-provider-pbs/commit/9b8729d564a19191c98f8bd500b92d372b42f975))
* Resolve VM integration test failures ([2ef8285](https://github.com/yavasura/terraform-provider-pbs/commit/2ef82857e133f407a234cecd2214c6d46daa3181))
* restore correct method signature in s3_providers.go ([768bc63](https://github.com/yavasura/terraform-provider-pbs/commit/768bc633b33ec7708218b47764a3588610b29cd9))
* run integration tests sequentially to avoid PBS resource contention ([66de76e](https://github.com/yavasura/terraform-provider-pbs/commit/66de76e1adf6c2de20b755c0729520be6cb13170))
* Separate metrics tests into subdirectories to avoid conflicts ([b024e05](https://github.com/yavasura/terraform-provider-pbs/commit/b024e0538381c28b22a03a2231d853c86329ddc4))
* Shorten job IDs to meet PBS 32 character limit ([58b1c00](https://github.com/yavasura/terraform-provider-pbs/commit/58b1c005ed18e36750a0b823cb4f4a307999d169))
* Simplify provisioning - use Debian's stock kernel, no reboot ([146ea15](https://github.com/yavasura/terraform-provider-pbs/commit/146ea15d7416869e9b560b0b8bee6a895de1d3c7))
* Split datasource tests into separate top-level directories ([6519ffe](https://github.com/yavasura/terraform-provider-pbs/commit/6519ffe301fd8beb8027d92e0f2d4f1c356b6745))
* Update AWS S3 tests to use us-east-1 region ([4815513](https://github.com/yavasura/terraform-provider-pbs/commit/48155134c44e414e72de72c888940608f3c4e9cc))
* Update HCL test provider version to match CI build ([51c4b36](https://github.com/yavasura/terraform-provider-pbs/commit/51c4b36d7722e9e954fd6765f0e29995be70bbea))
* Update remotes test to use environment variables correctly ([69a2b8f](https://github.com/yavasura/terraform-provider-pbs/commit/69a2b8f0aeb2ea2549a310d607ec75107eceac71))
* Update test paths in workflows and scripts ([58ad101](https://github.com/yavasura/terraform-provider-pbs/commit/58ad1017e293f31a5da3a467850b62b354bf9aaf))
* upgrade Terraform version in CI from 1.7.0 to 1.9.8 ([da4fa35](https://github.com/yavasura/terraform-provider-pbs/commit/da4fa3519771ac198e900fca7520863b0516b3bd))
* Use absolute paths for S3 datastore in immutability test ([330f445](https://github.com/yavasura/terraform-provider-pbs/commit/330f44551587c89caf37c4ded418fb9b7211c407))
* Use correct endpoint format for S3 datasource test ([bfa0673](https://github.com/yavasura/terraform-provider-pbs/commit/bfa0673dde19664b53cda89df9599c6039c18523))
* Use GH_PAT for cross-repository workflow dispatch ([c2961f0](https://github.com/yavasura/terraform-provider-pbs/commit/c2961f0b2fb4d2e3ad92c0d19d964a136cbe062e))
* Use unique test IDs and trigger cleanup workflow ([e406eeb](https://github.com/yavasura/terraform-provider-pbs/commit/e406eeb1155c3ffac75e8a39a9cbc74594157e60))
* Use workflow ID instead of filename for cleanup dispatch ([b71e245](https://github.com/yavasura/terraform-provider-pbs/commit/b71e24524c939ec825a04f9690a527672f2d24e3))

## Changelog

All notable changes to this project will be documented in this file.
