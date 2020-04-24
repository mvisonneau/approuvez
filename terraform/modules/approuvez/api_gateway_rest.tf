//
// REST - Used for Slack callbacks
//

resource "aws_api_gateway_rest_api" "rest" {
  name        = "ApprouvezREST"
  description = "REST entrypoint for approuvez's Slack callbacks"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "rest" {
  rest_api_id = aws_api_gateway_rest_api.rest.id
  parent_id   = aws_api_gateway_rest_api.rest.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "rest" {
  rest_api_id   = aws_api_gateway_rest_api.rest.id
  resource_id   = aws_api_gateway_resource.rest.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "rest" {
  rest_api_id = aws_api_gateway_rest_api.rest.id
  resource_id = aws_api_gateway_method.rest.resource_id
  http_method = aws_api_gateway_method.rest.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.slack_callback.invoke_arn
}

resource "aws_api_gateway_method" "rest_root" {
  rest_api_id   = aws_api_gateway_rest_api.rest.id
  resource_id   = aws_api_gateway_rest_api.rest.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "rest_root" {
  rest_api_id = aws_api_gateway_rest_api.rest.id
  resource_id = aws_api_gateway_method.rest_root.resource_id
  http_method = aws_api_gateway_method.rest_root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.slack_callback.invoke_arn
}

resource "aws_api_gateway_deployment" "rest" {
  depends_on = [
    aws_api_gateway_integration.rest,
  ]

  rest_api_id = aws_api_gateway_rest_api.rest.id
  stage_name  = local.stage_name
}
