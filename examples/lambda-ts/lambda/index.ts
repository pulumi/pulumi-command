import { Context } from "aws-lambda";

export async function handler(event: any, context: Context) {
  return context.functionName;
}
