package common

import "os"

const (
	awsLambdaEnvVar = "AWS_LAMBDA_FUNCTION_NAME"
)

func IsLambdaEnvironment() bool {
	return os.Getenv(awsLambdaEnvVar) != ""
}
