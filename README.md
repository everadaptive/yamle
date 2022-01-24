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

## Commands
encrypt with local keys
```
yamle --encrypt example/test_mixed_test_keys.yaml
```

decrypt with local keys
```
yamle --decrypt example/test_mixed_test_keys.yaml
```

generate cluster keys
```
yamle --cluster-key --gen-keys
```

encrypt with cluster keys
```
yamle --cluster-key --encrypt example/test_unencrypted.yaml > example/test_encrypted_cluster_keys.yaml
```

decrypt with cluster keys
```
yamle --cluster-key --decrypt example/test_encrypted_cluster_keys.yaml
```

## Building

```
go build ./...
```