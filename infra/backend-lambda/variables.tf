
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

variable "dlq_name" {
  type = string
  default = "omni-stack-dlq"
}

# Lambda configuration variables
variable "lambda_runtime" {
  type = string
  default = "java25"
  description = "Lambda runtime"
}

variable "micronaut_handler" {
  type = string
  default = "io.micronaut.function.aws.proxy.payload2.APIGatewayV2HTTPEventFunction"
  description = "Lambda function handler"
}

variable "omni_users_lambda_source_path" {
  type = string
  description = "Path to JAR file for omni-users lambda"
}

variable "omni_email_lambda_source_path" {
  type = string
  description = "Path to JAR file for omni-users lambda"
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

variable "db_table_arn" {
  type = string
  description = "ARN of DynamoDB table"
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
