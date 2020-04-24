// TODO: Make this one optional as this is working account wise
resource "aws_api_gateway_account" "default" {
  cloudwatch_role_arn = aws_iam_role.api_gateway.arn

  depends_on = [
    aws_iam_role_policy_attachment.api_gateway_cloudwatch_push_logs,
  ]
}

// create a role for API gateways
data "aws_iam_policy_document" "api_gateway_assume_role_policy" {
  statement {
    principals {
      type = "Service"

      identifiers = [
        "apigateway.amazonaws.com",
      ]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "api_gateway" {
  name               = "APIGatewayCloudwatchPushLogs"
  assume_role_policy = data.aws_iam_policy_document.api_gateway_assume_role_policy.json
}

// assign AmazonAPIGatewayPushToCloudWatchLogs policy
resource "aws_iam_role_policy_attachment" "api_gateway_cloudwatch_push_logs" {
  role       = aws_iam_role.api_gateway.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}
