output "rest_api_id" {
  value = aws_api_gateway_rest_api.rest.id
}

output "websocket_api_id" {
  value = aws_apigatewayv2_api.websocket.id
}

output "stage" {
  value = local.stage_name
}