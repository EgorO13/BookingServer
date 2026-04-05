# Room Booking Service

Сервис бронирования переговорок. Соответствует спецификации `api.yaml`.

## Запуск проекта

1. **Склонировать репозиторий**
   ```bash
   git clone https://github.com/internships-backend/test-backend-EgorO13.git
   cd task
2. **Запустить сервисы**
```bash
    make up
```    
Поднимает PostgreSQL и приложение через `docker-compose up -d`
3. **Наполнить базу тестовыми данными (опционально)**
```bash
    make seed
``` 
4. **Получить тестовый JWT**

Для роли user:
```bash
curl -X POST http://localhost:8080/dummyLogin \
-H "Content-Type: application/json" \
-d '{"role":"user"}'
```
Для роли admin:
```bash
curl -X POST http://localhost:8080/dummyLogin \
-H "Content-Type: application/json" \
-d '{ "role":"admin"}'
```
5. **Остановить работу**
```bash
    make down
```  

## Тестирование
- Юнит-тесты
```bash
    make unit-tests
```  
Появится файл `coverage.html`. Запустив его можно посмотреть покрытие
- E2E-тесты
  (Сервер должен быть запущен - `make up`)
```bash
    make e2e-tests
```

## Генерация слотов
Слоты генерируются при создании расписания на 30 дней вперёд. С запасом относительно 7 дней, но не очень много, что позволяет избежать избыточного разрастания таблицы `slots`.