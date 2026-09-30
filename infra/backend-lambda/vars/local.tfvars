environment = "local"
revision = "0.0.0-SNAPSHOT"
local_endpoint = "http://localhost:4566"
fundslock_lambda_source_path="../../backend/target/linux-amd64/fundslock-be"
enable_api_gw = true
api_spec_file = "openapi-spec.yaml"
lambda_env_vars = {
    "fundslock" = {
        "GRAPH_CONFIG_BASE64" = "W3siY2hhaW4iOiAic2Vwb2xpYSIsICJncmFwaFVybCI6ICJodHRwczovL2FwaS5zdHVkaW8udGhlZ3JhcGguY29tL3F1ZXJ5LzE3NTY3NzIvZnVuZHNsb2NrLXNlcG9saWEvdmVyc2lvbi9sYXRlc3QiLCAiYXBpS2V5IjoiNmMwMTE0YjEyOWJhYmZjOGM3M2QwY2E4ODU5MDNkNmEifSwgeyJjaGFpbiI6ICJhcmJpdHJ1bS1zZXBvbGlhIiwgImdyYXBoVXJsIjogImh0dHBzOi8vYXBpLnN0dWRpby50aGVncmFwaC5jb20vcXVlcnkvMTc1Njc3Mi9mdW5kc2xvY2stYXJiaXRydW0tc2Vwb2xpYS92ZXJzaW9uL2xhdGVzdCIsICJhcGlLZXkiOiAiNmMwMTE0YjEyOWJhYmZjOGM3M2QwY2E4ODU5MDNkNmEifV0=b"
        "JWT_SECRET_KEY" = "somekey"
    }
}