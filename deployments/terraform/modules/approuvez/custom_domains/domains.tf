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
  certificate_arn = aws_acm_certificate.all_domains.arn
}

resource "aws_apigatewayv2_domain_name" "rest" {
  domain_name = var.rest_domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.all_domains.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "slack" {
  api_id          = var.rest_api_id
  stage           = var.rest_api_stage
  domain_name     = aws_apigatewayv2_domain_name.rest.id
  api_key_mapping = "slack"
}

resource "aws_apigatewayv2_domain_name" "websocket" {
  domain_name = var.websocket_domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.all_domains.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "websocket" {
  api_id      = var.websocket_api_id
  stage       = var.websocket_api_stage
  domain_name = aws_apigatewayv2_domain_name.websocket.id
}
