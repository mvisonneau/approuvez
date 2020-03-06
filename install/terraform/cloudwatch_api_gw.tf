// TODO: Make this one optional as this is working account wise
resource "aws_api_gateway_account" "default" {
  cloudwatch_role_arn = aws_iam_role.api_gateway.arn

  depends_on = [
    aws_iam_role.api_gateway,
    aws_iam_role_policy_attachment.api_gateway_cloudwatch_push_logs,
  ]
}

// create a role for API gateways
resource "aws_iam_role" "api_gateway" {
  name = "APIGatewayCloudwatchPushLogs"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

// assign AmazonAPIGatewayPushToCloudWatchLogs policy
resource "aws_iam_role_policy_attachment" "api_gateway_cloudwatch_push_logs" {
  role       = aws_iam_role.api_gateway.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}
