// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package common

import "os"

const (
	awsLambdaEnvVar = "AWS_LAMBDA_FUNCTION_NAME"
)

func IsLambdaEnvironment() bool {
	return os.Getenv(awsLambdaEnvVar) != ""
}
