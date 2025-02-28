# Project Generation Orders

## Генератор заказов 0.3.0

### 1. Описание решения

Приложение - "Генератор заказов" реализует функционал по заказу товаров и услуг, например, для интернет-магазина. В нем реализована модульная система с возможностью репликации (Master-Replica) и партиционирования данных по городам в базе данных Postgres.

Исходные файлы приложения доступны в общем репозитории на GitHub: 
[GitHub Repository](https://github.com/nickfrolov2/orders-generator)

Локальный репозиторий для клонирования:
- SSH: `git@gitlaboff.devopsf.ru:mygroup1/test_project_frolovna.git`
- HTTPS: `https://gitlaboff.devopsf.ru/mygroup1/test_project_frolovna.git`

Подключение репозитория через Helm:

```bash
helm repo add orders https://raw.githubusercontent.com/nickfrolov2/orders-generator/refs/heads/main/
```

Запуск/установка производится с помощью команды:

```bash
helm install orders orders_generator/project_orders 
```

### 2. Описание компонентов и их взаимодействие

- **api-deployment.yml**  
  Приложение API `api-orders-generator` разворачивается из образа Docker контейнера `nikolayfrolov1986/api-orders-generator:2.1`, скачиваемого из публичного Docker Hub. Исходный файл с Dockerfile находится в корне проекта в папке GO.

- **api-service.yml**  
  Сервис для API.

- **ingress.yml**  
  Ingress, пробрасывает host и порт до api-service.

- **haproxy-deployment.yml, haproxy-service.yaml, haproxy-configmap.yml**  
  HAProxy выполняет роль балансировщика между двумя базами данных Postgres: `postgres-shards-0` и `postgres-shards-1`.

- **postgres-statefulset.yml, postgres-shards-service.yml, postgres-job1.yml, postgres-job2.yml, orders-replica-shards1-configmap.yml, orders-replica-shards2-configmap.yml**  
  В statefulset создаются два пода с контейнером Postgres: `postgres-shards-0` и `postgres-shards-1`. Для каждого пода монтируется `postgres-storage (PVC)`. Job1 подключается к `postgres-shards-0` и выполняет команду на создание таблицы `orders`, деля её на партиции по городам, создавая отдельную таблицу для каждого города, устанавливает уровень `logical` и создает публикацию для таблицы `orders`. Job2 работает аналогично для `postgres-shards-1`.

- **postgres-replica-statefulset.yml, orders-replica-configmap.yml, postgres-job3.yml**  
  Создается база Postgres `postgres-replica-0` для репликации данных с двух баз: `postgres-shards-0` и `postgres-shards-1`. Job3 создает таблицу `orders`, разделяет её на партиции по городам, создавая отдельную таблицу для каждого города, устанавливает уровень `replica` и создает подписку на таблицу `orders` баз `postgres-shards-0` и `postgres-shards-1`.

### 3. Разворачивание среды для приложения

Для работы приложения понадобятся инструменты:
- Helm
- Minikube (или аналогичные: Docker Compose, Rancher)

### 4. Запуск приложения

Подключение репозитория через Helm:

```bash
helm repo add orders https://raw.githubusercontent.com/nickfrolov2/orders-generator/refs/heads/main/
```

Запуск/установка производится с помощью команды:

```bash
helm install/upgrade orders orders_generator/project_orders 
```

Вы должны увидеть успешную установку:

```bash
PS C:\Users\1> helm upgrade orders orders_generator/project_orders

Release "orders" has been upgraded. Happy Helming!

NAME: orders
LAST DEPLOYED: Fri Feb 28 21:33:29 2025
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE: None

NOTES:

# Project Orders by Frolov N.A. (DevOps-kurs from RTK)

# Приложение "Генератор заказов" успешно установлено!

# Приятного использования! :-)
```

Для Windows/Linux/Mac:

В файле `hosts` добавляем строку:

```
<IP кластера> generator-orders.test
```

В браузере проверяем по ссылке: [http://generator-orders.test](http://generator-orders.test/)
```

### Примечание:
Убедитесь, что все компоненты настроены корректно и кластеры работают, прежде чем переходить к тестированию приложения.