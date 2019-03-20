# amazon-vpc-cni-k8s-tool
amazon-vpc-cni-k8s-tool for IP garbage collection

 ./cni-k8s-tool ip-gc --free-after=1s

Requires below IAM permissions:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "ec2:UnassignPrivateIpAddresses",
            "Resource": "*"
        }
    ]
}
```
