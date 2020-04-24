//
// CUSTOM DOMAIN
//
// Annoying to not be able to have both REST and WS API under the same domain name...
// TODO: Merge everything under the same domain once this become possible


resource "aws_acm_certificate" "all_domains" {
  domain_name       = var.rest_domain_name
  validation_method = "DNS"

  subject_alternative_names = [
    var.websocket_domain_name
  ]
}

resource "aws_acm_certificate_validation" "all_domains" {
  certificate_arn = aws_acm_certificate.approuvez.arn
}

resource "aws_apigatewayv2_domain_name" "rest" {
  domain_name = var.rest_domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.all_domains.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "example" {
  api_id      = aws_apigatewayv2_api.approuvez.id
  domain_name = aws_apigatewayv2_domain_name.rest.id
  stage       = aws_apigatewayv2_stage.approuvez_default.name
  base_path   = "slack"
}

resource "aws_api_gateway_base_path_mapping" "approuvez_slack" {
  api_id      = aws_api_gateway_rest_api.approuvez.id
  stage_name  = aws_api_gateway_deployment.approuvez.stage_name
  domain_name = aws_api_gateway_domain_name.approuvez_sph_re.domain_name
  base_path   = "slack"
}

resource "aws_apigatewayv2_domain_name" "websocket" {
  domain_name = "approuvez-ws.sph.re"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.approuvez_sph_re.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "websocket" {
  api_id      = aws_apigatewayv2_api.approuvez.id
  domain_name = aws_apigatewayv2_domain_name.approuvez_ws_sph_re.id
  stage       = aws_apigatewayv2_stage.approuvez_default.name
}
