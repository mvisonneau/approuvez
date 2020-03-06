resource "aws_api_gateway_rest_api" "approuver" {
  name        = "Approuver"
  description = "Entrypoint for approuver Slack callbacks"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "approuver" {
  rest_api_id = aws_api_gateway_rest_api.approuver.id
  parent_id   = aws_api_gateway_rest_api.approuver.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "approuver" {
  rest_api_id   = aws_api_gateway_rest_api.approuver.id
  resource_id   = aws_api_gateway_resource.approuver.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "approuver" {
  rest_api_id = aws_api_gateway_rest_api.approuver.id
  resource_id = aws_api_gateway_method.approuver.resource_id
  http_method = aws_api_gateway_method.approuver.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.slack_callback.invoke_arn
}

resource "aws_api_gateway_method" "approuver_root" {
  rest_api_id   = aws_api_gateway_rest_api.approuver.id
  resource_id   = aws_api_gateway_rest_api.approuver.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "approuver_root" {
  rest_api_id = aws_api_gateway_rest_api.approuver.id
  resource_id = aws_api_gateway_method.approuver_root.resource_id
  http_method = aws_api_gateway_method.approuver_root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.slack_callback.invoke_arn
}

resource "aws_api_gateway_deployment" "example" {
  depends_on = [
    aws_api_gateway_integration.approuver,
  ]

  rest_api_id = aws_api_gateway_rest_api.approuver.id
  stage_name  = "test"
}
