* Настраиваем го
```bash
mkdir ~/go
mkdir ~/go/src
echo "export GOPATH=~/go" >> ~/.profile # (или .bash_profile)
```

* Настраиваем IDE - включить поддержку модулей и выставить GOPATH в goland

* Настраиваем пакеты в го
```bash
go mod init mainmod
go mod tidy
```

```bash
docker build -t short-link-web:latest .
docker tag short-link-web:latest cr.yandex/crp2k2ebuohl073qpbfn/pupa-lupovich/short-web:latest
cat key.json | docker login --username json_key --password-stdin cr.yandex
docker push cr.yandex/crp2k2ebuohl073qpbfn/pupa-lupovich/short-web:latest
```