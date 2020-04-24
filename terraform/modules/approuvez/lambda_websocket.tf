// WEBSOCKETS

data "archive_file" "websocket" {
  type        = "zip"
  source_file = "${path.module}/websocket"
  output_path = "${path.module}/websocket.zip"
}

resource "aws_lambda_function" "websocket" {
  function_name    = "ApprouvezWebsocket"
  runtime          = "go1.x"
  filename         = "${path.module}/websocket.zip"
  role             = aws_iam_role.lambda.arn
  source_code_hash = data.archive_file.websocket.output_base64sha256
  handler          = "websocket"

  environment {
    variables = {
      API_GATEWAY_WEBSOCKET_ENDPOINT = local.api_gateway_websocket_endpoint
    }
  }
}

resource "aws_lambda_permission" "websocket_api_gateway" {
  count = length(local.websocket_routes)

  statement_id  = "AllowInvocationFromAPIGatewayForRoute${count.index}"
  action        = "lambda:InvokeFunction"
  function_name = "ApprouvezWebsocket"
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_apigatewayv2_api.websocket.execution_arn}/*/${local.websocket_routes[count.index]}"

  depends_on = [
    aws_lambda_function.websocket
  ]
}