# kubectl-track

[![asciicast](https://asciinema.org/a/i6UK3LE0WAZgWhg4o2HSupBim.svg)](https://asciinema.org/a/i6UK3LE0WAZgWhg4o2HSupBim)

`kubectl-track` monitors and displays changes for a specified Kubernetes resource, highlighting differences between resource versions to help in debugging and understanding resource evolution. It continuously tracks the resource, displaying each resource version on change, offering clear insights and facilitating troubleshooting. This makes it particularly useful when working with Kubernetes operators and reconciler loops.

## Installation

```shell
go install github.com/semihbkgr/kubectl-track@latest
```

## Usage

```shell
kubectl-track monitors and displays changes for a Kubernetes resource, highlighting differences between resource versions to help in debugging and understanding resource evolution

Usage:
  kubectl-track <resource> <name> [flags]

Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/semih/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for kubectl-track
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --log-file string                file to write logs to
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
  -v, --version                        version for kubectl-track
```
