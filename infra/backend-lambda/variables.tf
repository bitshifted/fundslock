
variable "environment" {
  description = "The environment for the stack. Used for tagging and other environment specific configuration."
  type        = string
}

variable "revision" {
  description = "The version of the stack. Used for tagging and other version specific configuration."
  type        = string
}

variable "aws_region" {
  type = string
  default = "eu-central-1"
}

variable "local_endpoint" {
  type = string
  default = null
}


# Lambda configuration variables
variable "lambda_runtime" {
  type = string
  default = "java25"
  description = "Lambda runtime"
}


variable "fundslock_lambda_source_path" {
  type        = string
  description = "Path to fundslock binary to be deployed as Lambda function"
}


variable "lambda_env_vars" {
  type = map(map(string))
  default = {}
  description = "Environment variables for lambda functions"
}

variable "lambda_encryption_key_alias" {
  type = string
  default = null
  description = "Alias KMS key for Lambda environment variable encryption"
}


# APi Gateway configuration variables
variable "enable_api_gw" {
  type = bool
  description = "Enable/disable API Gateway"
}

variable "api_spec_file" {
  type = string
  description = "Path to OpenAPI spec file"
}

variable "api_cors_origins" {
  type = list(string)
  description = "List of allowed origins for CORS"
}
