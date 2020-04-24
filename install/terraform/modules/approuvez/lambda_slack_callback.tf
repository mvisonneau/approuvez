data "archive_file" "slack_callback" {
  type        = "zip"
  source_file = "${path.module}/slack_callback"
  output_path = "${path.module}/slack_callback.zip"
}

resource "aws_lambda_function" "slack_callback" {
  function_name    = "ApprouvezSlackCallback"
  runtime          = "go1.x"
  filename         = "${path.module}/slack_callback.zip"
  role             = aws_iam_role.lambda.arn
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
  function_name = "ApprouvezSlackCallback"
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.rest.execution_arn}/*/*/*"

  depends_on = [
    aws_lambda_function.slack_callback,
  ]
}
