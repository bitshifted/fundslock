
data "aws_caller_identity" "current" {}

module "nosrv" {
  source = "git@github.com:bitshifted/cloud-tools.git//nosrv?ref=nosrv-1.1.0"

  # Lambda function definitions
  lambda_defs = {
    "fundslock" = {
      function_name = "fundslock"
      handler = "bootstrap"
      runtime = "provided.al2023"
      zip_archive_config = {
        source_file = "${var.fundslock_lambda_source_path}"
        output_dir = "/tmp/fundslock-lambda"
      }
      environment_variables = var.lambda_env_vars["fundslock"]

      execution_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
          {
            Sid = "AllowSSMParameterAccess"
            Effect = "Allow"
            Action = ["ssm:GetParameter"]
            Resource = [
              "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${var.environment}/fundslock/*"
            ]
          }
        ]
      })
    }
  }

  # Api Gateway configuration
  enable_api_gateway = var.enable_api_gw
  api_name = "fundslock-api"
  openapi_spec_path = "${path.cwd}/${var.api_spec_file}"
  api_cors_config = {
    allow_origins = var.api_cors_origins
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["Accept", "Authorization", "Content-Type", "X-CSRF-Token"]
    expose_headers = ["Link"]
    allow_credentials = true
  }

  environment = var.environment
  revision = var.revision

  tags = {
    "Application" = "fundslock"
  }
}


output "lambda_arns" {
  value = module.nosrv.lambda_function_arn
  description = "ARNs of lambda functions"
}
