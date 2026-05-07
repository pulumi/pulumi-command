pulumi {
  required_providers {
    aws = {
      source  = "pulumi/aws"
      version = ">= 7.0.0"
    }
  }
}

variable "aws_region" {
  type        = string
  description = "AWS region for the Lambda function."
}

resource "aws_iam_role" "lambda_role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_lambda_function" "lambda_function" {
  name    = "f"
  publish = true
  role    = aws_iam_role.lambda_role.arn
  handler = "index.handler"
  runtime = "nodejs20.x"
  code    = fileArchive("./handler")
}

resource "command_local_command" "invoke" {
  create = "aws lambda invoke --function-name \"$FN\" --payload '{\"stackName\": \"${pulumi.stack}\"}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '\"'  && rm out.txt"

  environment = {
    FN         = aws_lambda_function.lambda_function.arn
    AWS_REGION = var.aws_region
    AWS_PAGER  = ""
  }

  depends_on = [aws_lambda_function.lambda_function]
}

output "output" {
  value = command_local_command.invoke.stdout
}
