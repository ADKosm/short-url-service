- создаем кластер в Minikube
```
minikube delete # если до этого что-то было
minikube start --vm-driver=hyperkit
```

- запускаем под с echo сервером

```
kubectl run echo-server --image=ealen/echo-server
kubectl get pods -o yaml | grep IP # достаем IP пода
```

- запускаем еще один под с убунтой

```
kubectl run -it ubuntu --image=ubuntu:20.04
apt update && apt install curl
curl 127.0.0.3/hello
```

---

- переписываем то же самое через yaml + добавляем сервис
```
kubectl apply -f step00-echo-service.yaml
```

- в ubuntu-поде видим:
  - наличие переменной окружения
  - резолвинг по dns

---

- запускаем в kubernetes монгу:

```
kubectl run my-mongo-pod --image=mongo:4.4 # запустить под руками
kubectl get pods my-mongo -o yaml
kubectl apply -f step00-run-mongo-in-kubernetes.yaml # то же самое, тольк через deployment и yaml
```

- подключаемся к контейнеру в поде 

```
kubectl get pods
kubectl exec -it dummy-mongo-78d6555548-fd4m4 bash
```

- смотрим адрес пода в 

- регистрируемся/логинимся на [hub.docker.com](https://hub.docker.com) и создаем там публичный репозиторий

- логинимся в docker cli:
```
docker login
```

- собираем docker-образ сервиса и пушим его в Central Docker Registry

```
docker build src -t lfyuomrgylo/url-shortener:v1
docker push lfyuomrgylo/url-shortener:v1
```

- создаем под, в котором запущен наш сервис вместе с монгой и смотрим на них

```
kubectl create -f step0-service-and-mongo-in-one-pod.yam
kubectl get replicasets
kubectl get pods --selector service-name=url-shortener -o yaml
```
