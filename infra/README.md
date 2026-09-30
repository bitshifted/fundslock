
# Run Floci in Docker container

Floci provides local AWS environment and services for testing. To run it, use:

```
docker compose up -f docker-compose-local.yaml up
```

# Configure AWS CLI for local connection

Create new profile for local development:

```
aws configure --profile local
```

Enter mock data for required information:

```
AWS Access Key ID [None]: mock
AWS Secret Access Key [None]: mock
Default region name [None]: eu-central-1
Default output format [None]: JSON
```

Connect to Floci using command line:

```
aws --endpoint-url=http://localhost:4566 s3 ls --profile local
```

## Configure endpoint within profile

To avoid having to type `--endpoint` parameter each time you run the command, you can set it inside profile. Open `~/.aws/config` file and add the following to `local` profile section:

```
[profile local]
....
endpoint_url = http://localhost:4566
```

Now command becomes:

```
aws --profile local s3 ls
```

## Set profile as environment variable

To avoid having to specify profile, set it as environment variable:

```
export AWS_PROFILE=local
```

Now the command becomes: `aws s3 ls`
