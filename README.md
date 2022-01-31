# YAMLE - Easy YAML Encryption

## YAML Format
```YAML
# A value to be encrypted
value1: !encrypt "I'm a secret"

# A value that is already encrypted
value2: !encrypted FYooOnahFLnyoZwrJ09X+ppKbSc6uLD2/uIdPXCuElPmDBcTJeSCR4YcTDwcMbym79bfqkebnsh3rbw0oVjReRre8L7uuHHOEYscCXd+dpV5/ipBWcdSgxQkniHBXYSy7ilYLDB2PIC54/e/Zj5Jbkzq10shRqNoZyDW7xSIOo/5DCtdLawXdNKEChj8zFsqyv85AKHGpW6kHM45buuivxExFS6Tm0lR38CHGIVJWtTqu3ITYxT3fCqsvcPg1DMUUkR8Nus8uM18PLwsZ1bsEnr0ioOZfIufWOPj+XI6JQmusNmWOnWGmKQ14Q1dFKfRsQb62SXPA4zaJIvUxjYawQ==

# A value that won't be touched
value3: "I'm not a secret"
```

## Usage
```
A simple YAML encrypter

Usage:
  yamle [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  decrypt     Decrypt a YAML file
  encrypt     Encrypt a YAML file
  gen-keys    Generate RSA keys to use for encryption and decryption
  help        Help about any command
  version     Print the version number of yamle

Flags:
      --cluster-key                    use keys stored in Kubernetes
      --cluster-key-name string        the name of the secret in Kubernetes
      --cluster-key-namespace string   the namespace of the secret in Kubernetes
      --config string                  config file (default is ./.yamle.yaml)
  -h, --help                           help for yamle
      --in-place                       encrypt or decrypt file in place
      --key-file string                path to a PEM formatted key
      --key-size int                   size of generated keys (default is 2048) (default 2048)
```

## Commands
encrypt with local keys
```
yamle encrypt example/test_mixed_test_keys.yaml
```

decrypt with local keys
```
yamle decrypt example/test_mixed_test_keys.yaml
```

generate cluster keys
```
yamle gen-keys --cluster-key --cluster-key-name=yamle-keys --cluster-key-namespace=default
```

encrypt with cluster keys
```
yamle encrypt--cluster-key example/test_unencrypted.yaml > example/test_encrypted_cluster_keys.yaml
```

decrypt with cluster keys
```
yamle decrypt --cluster-key example/test_encrypted_cluster_keys.yaml
```

## Building

```
go build ./...
```