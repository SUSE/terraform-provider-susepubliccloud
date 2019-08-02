[![Build Status](https://travis-ci.org/SUSE/terraform-provider-susepubliccloud.svg?branch=master)](https://travis-ci.org/SUSE/terraform-provider-susepubliccloud)

The purpose of this project is to define a set of [data sources](https://www.terraform.io/docs/configuration/data-sources.html)
that make it easier to find the resources managed by SUSE on the different public clouds.

The provider gathers data by accessing the public instance of
[public-cloud-info-service](https://github.com/SUSE-Enceladus/public-cloud-info-service)
managed by SUSE.

You can find more information about the resources managed by SUSE on the public clouds
here:

* [Blog post](https://www.suse.com/c/riddle-me-this/)
* [Offical cli tool](https://github.com/SUSE-Enceladus/public-cloud-info-client)

## Using the Provider

### Data source `susepubliccloud_image_ids`

Use this data source to get a list of image IDs matching
the specified criteria.

Example use:

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

#### Argument reference

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

#### Attributes reference

`ids` is set to the list of images IDs, sorted by publication time according to
`sort_ascending`.

## Installing the Provider

openSUSE and SUSE packages are being built inside of the
[systemsmanagement:terraform terraform-provider-susepubliccloud](https://build.opensuse.org/package/show/systemsmanagement:terraform/terraform-provider-susepubliccloud)
project on the [Open Build Service](https://build.opensuse.org/).

The packages can be installed by visiting the [dedicated page](https://software.opensuse.org/package/terraform-provider-susepubliccloud?search_term=terraform-provider-susepubliccloud)
on [software.opensuse.org](https://software.opensuse.org).

## Developing the Provider

If you wish to work on the provider, you'll need:

* [Terraform](https://www.terraform.io/downloads.html) 0.10+
* [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:SUSE/terraform-provider-susepubliccloud
...
```

To compile the provider execute the following command:

```sh
$ make build
```

This will produce the `terraform-provider-susepubliccloud` binary that can be
copied to your `~/.terraform.d/plugins/<GOOS>_<GOARCH>` directory.

More instructions about how to instead use a custom-built provider in your
Terraform environment can be found
[here](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).

### Testing

Unit tests can be run via:

```sh
$ make test
```

Unit test coverage can be seen by executing:

```sh
$ make test-coverage
```
