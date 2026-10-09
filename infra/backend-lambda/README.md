
For local development, export mock AWS credentials and set local profile as current. See [here](../README.md)

# Initialize Terraform

```
terraform init -backend-config=./vars/local.s3.tfbackend
```

## Reinitialize for differennt backend

Running the same Terraform plan for different backend requires reconfiguring backend.

```
terraform init -backend-config=./vars/local.s3.tfbackend -reconfigure
```

# Run plan

Export dynamic variables:
```
export TF_VAR_revision=0.0.1
```

```
terraform plan -var-file=./vars/local.tfvars
```

