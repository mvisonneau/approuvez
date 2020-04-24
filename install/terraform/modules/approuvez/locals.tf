locals {
  stage_name = "default"

  websocket_routes = [
    "$connect",
    "$disconnect",
    "$default",
    "get_connection_id",
  ]

  api_gateway_websocket_endpoint = "${replace(aws_apigatewayv2_api.websocket.api_endpoint, "wss://", "")}/${local.stage_name}"
}