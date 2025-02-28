# Используем официальный образ Go для сборки
FROM golang:1.23.6 AS builder  

# Создаем рабочий каталог
WORKDIR /app  

# Копируем файлы проекта
COPY . .

# Проверяем зависимость модуля
RUN go mod tidy  

# Составляем исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp 

# Новый минимальный контейнер:
FROM alpine:latest  
WORKDIR /root/  

# Копируем скомпилированный бинарник в новый контейнер
COPY --from=builder /app/myapp /root/myapp 
COPY static static

# Запускаем сервер
CMD ["./myapp"]