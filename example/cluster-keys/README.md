# Using RBAC

With cluster keys, you can use RBAC to enable users to encrypt values while not
being able to decrypt them.

encrypt with cluster keys
```
yamle --cluster-key --key-secret test-public --encrypt example/test_unencrypted.yaml > example/test_encrypted_cluster_keys.yaml
```

decrypt with cluster keys
```
yamle --cluster-key --decrypt example/test_encrypted_cluster_keys.yaml
```