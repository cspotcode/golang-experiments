## Purpose

* Docker analogue to `git credential store` for directly storing credentials, delegating the the appropriate credential
store or storing in .docker/config.json

* Rapidly authenticate against ECR repositories in many AWS regions.

* Avoid slowness of `docker login`, which verifies the credentials.  When authenticating for ECR, we already trust that
the AWS SDK is giving us valid credentials.

The following will store credentials for ECR in the given account and regions, assuming you have configured AWS CLI
credentials.  For example, having `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` will work.

```
docker-credential store aws --regions us-west-2,us-east-1 <aws account id>
```

## Usage

Install from source:

```
go install github.com/cspotcode/golang-experiments/docker-credential
```

Or download a prebuilt binary from the most recent "Github Actions" build: https://github.com/cspotcode/golang-experiments/actions

Get usage info:

```
docker-credential --help
```

Run like this:

```
docker-credential --verbose store aws --regions us-east-1,us-west-2 <your AWS account number>
```

For convenience, I have this function in my shell profile which logs into all my team's regions:
```
function ecrlogin {
    docker-credential --verbose store aws --regions ap-northeast-1,ap-northeast-2,ap-south-1,ap-southeast-1,ap-southeast-2,ca-central-1,eu-central-1,eu-north-1,eu-west-1,eu-west-2,eu-west-3,sa-east-1,us-east-1,us-east-2,us-west-1,us-west-2 <my AWS account number>
}
```

To login:

```
ecrlogin
```

## Why not amazon-ecr-credential-helper?

It causes unrelated commands to slow down: https://github.com/moby/moby/issues/31517

For example, `docker build` appears to hang because it invokes the ECR credential helper for all configured repositories.
It does this even when it's completely unnecessary, because it wants to send
all known credentials to the docker engine before the build starts.  When your shell does not have access to AWS
credentials, the helper is even slower.

## How to use amazon-ecr-credential-helper

If you still want to use [amazon-ecr-credential-helper](https://github.com/awslabs/amazon-ecr-credential-helper), install it:

```
go get -u github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login
```

...then set it up in my `~/.docker./config.json`:

```
    // Add this to ~/.docker/config.json, replacing with the relevant aws account id
  "credHelpers": {
    "<aws account id>.dkr.ecr.ap-east-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ap-northeast-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ap-northeast-2.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ap-south-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ap-southeast-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ap-southeast-2.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.ca-central-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.eu-central-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.eu-north-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.eu-west-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.eu-west-2.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.eu-west-3.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.me-south-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.sa-east-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.us-east-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.us-east-2.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.us-west-1.amazonaws.com": "ecr-login",
    "<aws account id>.dkr.ecr.us-west-2.amazonaws.com": "ecr-login"
  },
```

## TODO

* figure out identitytoken in struct; are we supposed to be converting username and password into identity token?
