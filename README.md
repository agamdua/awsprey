AWSprey - Bird is the word
==========================

### Installation

Two parts:

```
$ brew tap agamdua/homebrew-awsprey
$ brew install awsprey
```

Set AWS credentials with priveleges to describe EC2 instances. This currently
uses your default AWS profile. Future versions will allow you to use
custom profiles.

Check if you have a file `~/.aws/credentials` that looks like:

```
[default]
aws_access_key_id = xxxxxxxxxxxxxx
aws_secret_access_key = xxxxxxxxxxxxxxx
```

If not, create one and add your values there.

### Usage

Once you are done setting up, assuming your EC2 instances are tagged
with `service` and `environment` tags:

```
$ awsprey list <service>:<environment>
```

### Example

If you have a service called `web` and want to check environment `staging`:

```
$ awsprey list web:staging
```


### Development

Run tests with

```
make test
```
