#Project Generation Orders

#Генератор заказов 0.3.0

<b>1. Описание решения</b>

Приложение - "Генератор заказов" реализует функционал по заказу товаров и услуг, например, для интернет магазина. Реализована модульная система с возможностью реплицирования (Master-Replica) и партицирования данных по городам в БД Postgres.

Исходные файлы приложения доступны в общем репозитории на github: https://github.com/nickfrolov2/orders-generator

В локальном репозитории для клонирования: git@gitlaboff.devopsf.ru:mygroup1/test_project_frolovna.git и https://gitlaboff.devopsf.ru/mygroup1/test_project_frolovna.git

Подключение репозитория через helm:
helm repo add orders https://raw.githubusercontent.com/nickfrolov2/orders-generator/refs/heads/main/
Запуск/установка производится с помощью команды:
helm install orders orders_generator/project_orders 


<b>2. Описание компонентов и их взаимподействие</b>

<b>api-deployment.yml</b> - приложение api api-orders-generator разворачивается из образа докер контейнера nikolayfrolov1986/api-orders-generator:2.1, происходит скачивание из публичного docker hub. Исходный файл с Dockerfile находится в корне проекта в папке GO/
<b>api-service.yml</b> - сервис для api
<b>ingress.yml</b> - ingress, пробрасывает host и порт до api-service
<b>haproxy-deployment.yml, haproxy-service.yaml, haproxy-configmap.yml</b> - haproxy выполняет роль балансировщика между двумя базами данных Postgres: 
postgres-shards-0
postgres-shards-1
<b>postgres-statefulset.yml, postgres-shards-service.yml, postgres-job1.yml, postgres-job2.yml, orders-replica-shards1-configmap.yml, orders-replica-shards2-configmap.yml <b> - В statefulset создается два пода с контейнером Postgres postgres-shards-0 и postgres-shards-1. Для каждого монтируется postgres-storage(pvc). Job1 подключается к postgres-shards-0 и запускает команду на создание таблицы orders, разделяет на партиции по городам, создает отдельную таблицу для каждого города, устанавливает уровень logical и создает публикацию таблицы orders. Job2 – работает аналогично Job1, только для postgres-shards-1.
<b>postgres-replica-statefulset.yml, orders-replica-configmap.yml, postgres-job3.yml</b> - создается база Postgres postgres-replica-0 для реплицирования данных с двух баз postgres-shards-0 и postgres-shards-1. Job3 создает таблицу orders, разделяет на партиции по городам, создает отдельную таблицу для каждого города, устанавливает уровень replica и создает подписку на таблицу orders баз postgres-shards-0 и postgres-shards-1.

<b>3. Разворачивание среды для приложения</b>
Для работы приложения понадобятся инструменты: helm, minikube (аналог, docker compose, rancher) 

<b>4. Запуск приложения.</b>
Подключение репозитория через helm:
helm repo add orders https://raw.githubusercontent.com/nickfrolov2/orders-generator/refs/heads/main/
Запуск/установка производится с помощью команды:
helm install/upgrade orders orders_generator/project_orders 
Должны увидеть успешную установку:
PS C:\Users\1> helm upgrade orders orders_generator/project_orders
Release "orders" has been upgraded. Happy Helming!
NAME: orders
LAST DEPLOYED: Fri Feb 28 21:33:29 2025
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE: None
NOTES:
#Project Orders by Frolov N.A.(DevOps-kurs from RTK)
#Приложение "Генератор заказов" успешно установлено!
#Приятного использования! :-)

Для windows/linux/mac:
В файле hosts добавляем строки -  “IP кластера”   generator-orders.test
В браузере проверяем по ссылке http://generator-orders.test/
