# ec2id
A CLI tool that retrieve the EC2 instance ID of specified instance's Name tag.

## Usage
Configure your AWS credentials, and then run as follows

```bash
# retrieve the latest launched instance id
$ ec2id
i-0691a69ff0914bae1
```

You can specify the instnace's Name tag
```
# retrieve the instance id with tag:Name = "sample"
$ ec2id sample
i-0acd9f178c934caea
```

