environment = "local"
revision = "0.0.0-SNAPSHOT"
local_endpoint = "http://localhost:4566"
fundslock_lambda_source_path="../../backend/target/linux-amd64/fundslock-be"
enable_api_gw = true
api_spec_file = "openapi-spec.yaml"
api_cors_origins = ["http://localhost:5173"]
lambda_env_vars = {
    "fundslock" = {
        "ACCESS_TOKEN_DURATION" = "900"
        "GRAPH_CONFIG_BASE64" = "3siY2hhaW4iOiJzZXBvbGlhIiwiZ3JhcGhVcmwiOiJodHRwczovL3NvbWVob3N0LmNvbS9ncmFwaCIsImFwaUtleSI6ImFiY2RlZiJ9XQ=="
        "JWT_SECRET_KEY" = "somekey"
    }
}