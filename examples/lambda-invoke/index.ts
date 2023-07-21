import { local } from "@pulumi/command";
import * as aws from "@pulumi/aws";

const f = new aws.lambda.CallbackFunction("f", {
    publish: true,
    callback: async (ev: any) => {
        const crypto = require("crypto");
        return crypto.randomBytes(ev.len/2).toString('hex');
    }
});

const rand = new local.Command("execf", {
    create: `aws lambda invoke --function-name "$FN" --payload '{"len": 10}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '"'  && rm out.txt`,
    environment: {
        FN: f.qualifiedArn,
        AWS_REGION: aws.config.region!,
        AWS_PAGER: "",
    },
})

export const output = rand.stdout;
