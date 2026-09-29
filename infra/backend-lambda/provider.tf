
provider "aws" {
  region     = var.aws_region

  # Set the endpoint to the LocalStack API
dynamic "endpoints" {
    for_each = var.local_endpoint != null ? [1] : []
  content {
    apigateway     = var.local_endpoint
    dynamodb = var.local_endpoint
    lambda = var.local_endpoint
  }
}

}