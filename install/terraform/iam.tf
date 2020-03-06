//
// LAMBDA
//

data "aws_caller_identity" "current" {}

// create a role for our lambda function
resource "aws_iam_role" "approuver_lambda_function" {
  name = "ApprouverLambdaFunction"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

// assign default execution permissions to make our function executable
resource "aws_iam_role_policy_attachment" "approuver_lambda_function_basic_execution_role" {
  role       = aws_iam_role.approuver_lambda_function.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// authorise access onto apigateway
resource "aws_iam_policy" "approuver_lambda_function_api_gateway_websocket" {
  name        = "ApprouverLambdaFunctionAPIGatewayWebsocket"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "execute-api:*",
      "Resource": "${local.api_gateway_websocket_arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "approuver_lambda_function_api_gateway_websocket" {
  role       = aws_iam_role.approuver_lambda_function.name
  policy_arn = aws_iam_policy.approuver_lambda_function_api_gateway_websocket.arn
}