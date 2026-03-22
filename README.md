# ood-emr-adapter

OOD compute adapter for [Amazon EMR Serverless](https://docs.aws.amazon.com/emr/latest/EMR-Serverless-UserGuide/getting-started.html). Translates Open OnDemand job lifecycle calls (submit / status / delete / info) to the EMR Serverless API.

## Commands

| Command | Description |
|---------|-------------|
| `submit` | Read a JSON job spec from stdin and submit it as an EMR Serverless job run. Prints `<application-id>/<job-run-id>` on success. |
| `status <id>` | Return OOD-normalised job status as JSON (`queued`, `running`, `completed`, `failed`, `cancelled`, `undetermined`). |
| `delete <id>` | Cancel a running EMR Serverless job run. |
| `info <id>` | Print the full `GetJobRun` API response as JSON. |

`<id>` must be in the form `<application-id>/<job-run-id>`, or just `<job-run-id>` when `--application-id` is set.

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--region` | `us-east-1` | AWS region |
| `--application-id` | _(empty)_ | EMR Serverless application ID (can also be set per-job in the spec) |

## Job spec (stdin for `submit`)

```json
{
  "application_id": "00abc123def456gh",
  "execution_role_arn": "arn:aws:iam::123456789012:role/EMRServerlessRole",
  "entry_point": "s3://my-bucket/scripts/wordcount.py",
  "entry_point_args": ["s3://my-bucket/input/", "s3://my-bucket/output/"],
  "spark_submit_parameters": "--conf spark.executor.cores=4 --conf spark.executor.memory=8g",
  "job_name": "wordcount-run-42",
  "env": {
    "MY_ENV_VAR": "value"
  }
}
```

## Open OnDemand cluster YAML example

```yaml
# config/clusters.d/emr-serverless.yml
v2:
  metadata:
    title: "Amazon EMR Serverless"
  login:
    host: "emr-serverless.internal"
  job:
    adapter: "script"
    submit: "/usr/local/bin/ood-emr-adapter submit --region us-east-1 --application-id 00abc123def456gh"
    status: "/usr/local/bin/ood-emr-adapter status --region us-east-1"
    delete: "/usr/local/bin/ood-emr-adapter delete --region us-east-1"
    info:   "/usr/local/bin/ood-emr-adapter info   --region us-east-1"
```

## Prerequisites

- Go 1.26+
- AWS credentials configured (environment variables, `~/.aws/credentials`, or IAM role)
- An existing EMR Serverless application

## Build

```bash
go build -o ood-emr-adapter .
```

## License

MIT — see [LICENSE](LICENSE).
