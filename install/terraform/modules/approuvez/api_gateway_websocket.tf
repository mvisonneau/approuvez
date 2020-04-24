//
// WEBSOCKET - Used for client connectivity
//


resource "aws_apigatewayv2_api" "websocket" {
  name        = "ApprouvezWebSocket"
  description = "Websocket entrypoint for approuvez's clients"

  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.route"
}

resource "aws_apigatewayv2_stage" "websocket_default" {
  api_id      = aws_apigatewayv2_api.websocket.id
  auto_deploy = true
  name        = local.stage_name
}

resource "aws_apigatewayv2_integration" "websocket" {
  api_id                    = aws_apigatewayv2_api.websocket.id
  integration_type          = "AWS_PROXY"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.websocket.invoke_arn
  content_handling_strategy = "CONVERT_TO_TEXT"
}

// Configure routes
resource "aws_apigatewayv2_route" "websocket" {
  count = length(local.websocket_routes)

  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = local.websocket_routes[count.index]
  target    = "integrations/${aws_apigatewayv2_integration.websocket.id}"
}