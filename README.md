# amazon-vpc-cni-k8s-tool
amazon-vpc-cni-k8s-tool for IP garbage collection

## Prerequisites

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

## Usage

 ./cni-k8s-tool list-free-IPs

 ./cni-k8s-tool registry-list

 ./cni-k8s-tool ip-gc --free-after=1s

## K8s native installation

A Dockerfile and K8s CronJob kubeyaml manifest is provided.

Build the Dockerfile and push to an image registry, then update the 'image:' in cni-k8s-tool_ip-gc_cron.yaml

Applying cni-k8s-tool_ip-gc_cron.yaml will create a CronJob in the kube-system namespace that runs the `./cni-k8s-tool ip-gc --free-after=1s` from above every minute.

Note: 'Completed' Jobs have been observed to pile up in the `kubectl get pods -n kube-system` output. Tools such as [onfido/k8s-cleanup](https://github.com/onfido/k8s-cleanup ) may also be desirable to clean those up.
