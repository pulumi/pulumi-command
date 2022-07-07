// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

import * as aws from "@pulumi/aws";
import { local } from "@pulumi/command";

export = async () => {
  const build = await local.run({
    command: "yarn && yarn tsc",
    dir: "lambda",
    archivePaths: ["*.js", "*.js.map"],
  });

  // Lambda
  const role = new aws.iam.Role("lambda-role", {
    assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({
      Service: "lambda.amazonaws.com",
    }),
  });
  const lambda = new aws.lambda.Function("ts-lambda", {
    role: role.arn,
    code: build.archive,
    handler: "index.handler",
    runtime: aws.lambda.Runtime.NodeJS16dX,
  });

  // Exports
  return {
    invokeArn: lambda.invokeArn,
    lambdaArchive: build.archive,
  };
};
