## How to encrypt your content


#### Encrypt plantext

```bash
# Encrypt plaintext for namespace 'test', secret name 'hello' and env is 'develop'
$ echo 'EnterYourPlaintextHere' | base64 | xargs curl 'http://localhost:8000/encrypt?namespace=test&name=hello&env=develop' -v -XPOST -d
```

#### Encrypt a file

```bash
# Encrypt file content for namespace 'test', secret name 'hello' and env is 'develop'
$ cat google-service-account.json | base64 | xargs curl 'http://localhost:8000/encrypt?namespace=test&name=hello&env=develop' -v -XPOST -d
```
