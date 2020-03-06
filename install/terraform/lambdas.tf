data "archive_file" "slack_callback" {
  type        = "zip"
  source_file = "${path.module}/slack_callback"
  output_path = "${path.module}/slack_callback.zip"
}

resource "aws_lambda_function" "slack_callback" {
  function_name    = "slack_callback"
  runtime          = "go1.x"
  filename         = "slack_callback.zip"
  role             = aws_iam_role.approuver_lambda_function.arn
  source_code_hash = data.archive_file.slack_callback.output_base64sha256
  handler          = "slack_callback"

  environment {
    variables = {
      API_GATEWAY_WEBSOCKET_ENDPOINT = local.api_gateway_websocket_endpoint
    }
  }
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowInvocationFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "slack_callback"
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.approuver.execution_arn}/*/*/*"
}

// WEBSOCKETS

data "archive_file" "websocket" {
  type        = "zip"
  source_file = "${path.module}/websocket"
  output_path = "${path.module}/websocket.zip"
}

resource "aws_lambda_function" "websocket" {
  function_name    = "websocket"
  runtime          = "go1.x"
  filename         = "websocket.zip"
  role             = aws_iam_role.approuver_lambda_function.arn
  source_code_hash = data.archive_file.websocket.output_base64sha256
  handler          = "websocket"

  environment {
    variables = {
      API_GATEWAY_WEBSOCKET_ENDPOINT = local.api_gateway_websocket_endpoint
    }
  }
}
