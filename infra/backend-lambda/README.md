
For local development, export mock AWS credentials and set local profile as current. See [here](../README.md)

# Initialize Terraform

```
terraform init -backend-config=./vars/local.s3.tfbackend
```

# Run plan

```
terraform plan -var-file=./vars/local.tfvars
```

