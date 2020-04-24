//
// LAMBDA
//

// Create a role for the lambda functions
data "aws_iam_policy_document" "lambda_assume_role_policy" {
  statement {
    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda" {
  name               = "ApprouvezLambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

// Assign default execution permissions to make the function executable
resource "aws_iam_role_policy_attachment" "lambda_basic_execution_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Add some IAM permissions to the functions
data "aws_iam_policy_document" "lambda_default_policy" {
  // Authorize the function to trigger the Websocket one
  // (this is only necessary for the slack-callback lambda function)
  statement {
    actions = [
      "execute-api:*",
    ]

    resources = [
      "${aws_apigatewayv2_api.websocket.execution_arn}/*",
    ]
  }
}

resource "aws_iam_role_policy" "lambda_default" {
  name   = "default"
  role   = aws_iam_role.lambda.name
  policy = data.aws_iam_policy_document.lambda_default_policy.json
}
