
module "nosrv" {
  source = "git@github.com:bitshifted/cloud-tools.git//nosrv?ref=nosrv-1.3.0"

  dlq_name = var.dlq_name
  sqs_defs = {
    "email_queue" = {
      name = "omni-email-queue"
    }
  }

  # Lambda function definitions
  lambda_defs = {
    "omni_users" = {
      function_name = "omni-users"
      handler = var.micronaut_handler
      runtime = var.lambda_runtime
      zip_archive_config = {
        source_file = "${var.omni_users_lambda_source_path}"
        output_dir = "/tmp/omnistack-lambda"
      }
      
      execution_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
          {
            Sid = "AllowDynamoDBAccess"
            Effect = "Allow"
            Action = [
              "dynamodb:GetItem",
              "dynamodb:PutItem",
              "dynamodb:UpdateItem",
              "dynamodb:DeleteItem",
              "dynamodb:Query",
              "dynamodb:Scan",
              "dynamodb:BatchGetItem",
              "dynamodb:BatchWriteItem"
            ]
            Resource = [
              "${var.db_table_arn}",
              "${var.db_table_arn}/index/*"
            ]
          },
          {
            Sid = "AllowSQSWriteAccess"
            Effect = "Allow"
            Action = [
              "sqs:SendMessage",
              "sqs:GetQueueAttributes",
              "sqs:GetQueueUrl"
            ]
            Resource = [
              module.nosrv.sqs_arns["email_queue"]
            ]
          }
          
        ]
      })
      environment_variables = var.lambda_env_vars["omni_users"]
      encryption_key_alias = var.lambda_encryption_key_alias
    }
    "omni_email" = {
      function_name = "omni-email"
      handler = var.micronaut_handler
      runtime = var.lambda_runtime
      zip_archive_config = {
        source_file = var.omni_email_lambda_source_path
        output_dir = "/tmp/omnistack-lambda"
      }
      
      execution_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
          {
            Sid = "AllowSQSReadAccess"
            Effect = "Allow"
            Action = [
              "sqs:ReceiveMessage",
              "sqs:DeleteMessage",
              "sqs:GetQueueAttributes",
              "sqs:ChangeMessageVisibility"
            ]
            Resource = [
              module.nosrv.sqs_arns["email_queue"]
            ]
          }
        ]
      })
      environment_variables = var.lambda_env_vars["omni_email"]
      encryption_key_alias = var.lambda_encryption_key_alias
    }
  }

  # Api Gateway configuration
  enable_api_gateway = var.enable_api_gw
  api_name = "omnistack-api"
  openapi_spec_path = "${path.cwd}/${var.api_spec_file}"

  environment = var.environment
  revision = var.revision

  tags = {
    "Application" = "omni-stack"
  }
}

output "sqs_arns" {
  value = module.nosrv.sqs_arns
  description = "ARNs of configured queue"
}

output "sqs_urls" {
  value = module.nosrv.sqs_urls
  description = "URLs of configured queue"
}

output "lambda_arns" {
  value = module.nosrv.lambda_function_arn
  description = "ARNs of lambda functions"
}
