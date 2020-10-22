# susepubliccloud_image_ids Data Source

Use this data source to get a list of image IDs matching
the specified criteria.

## Example Usage

```hcl
data "susepubliccloud_image_ids" "sles" {
  cloud      = "amazon"
  region     = "eu-central-1"
  state      = "active"
  name_regex = "suse-sles-15-sp1-byos.*-hvm-ssd-x86_64"
}

resource "aws_instance" "control_plane" {
  ami = "${data.susepubliccloud_image_ids.sles.ids[0]}"
  ...
}
```

### Argument Reference

* `cloud` - (Required) Name of the target cloud to use. Valid values: `amazon`,
  `google`, `microsoft` and `oracle`.
* `region` - (Required) One of the known regions in the cloud framework. Use the
  region identifiers as the provider describes them, for example `us-east-1` in
  Amazon EC2, or `East US 2` in Microsoft Azure.
* `state` - (Defaults to `active`) State of the image. Valid values:
  `active`, `inactive`, `deprecated`. Note well: the `deleted` state isn't
  accepted by the data source because these images would not be usable by
  terraform.
* `name_regex` - (Optional) A regex string to apply to the images list returned
  by the remote API managed by SUSE.
* `sort_ascending` - (Defaults to `false`) Used to sort by publication time.

**Note well:** the values accepted by `cloud`, `region` and `state` are the ones
specified [here](https://github.com/SUSE-Enceladus/public-cloud-info-service#server-design).

### Attributes Reference

* `ids` is set to the list of images IDs, sorted by publication time according to
`sort_ascending`.
