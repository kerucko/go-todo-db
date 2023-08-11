# todo_db

Туду-лист в виде сайта с использованием базы данных MySQL

## Инструкция по запуску:
1. Созать БД MySQL с двумы таблицами:
```sql
CREATE TABLE tasks (
  id INT NOT NULL AUTO_INCREMENT, 
  name VARCHAR(50) NOT NULL, 
  comment VARCHAR(100), 
  createDate DATE NOT NULL, 
  deadline DATE, 
  appointmentDate DATE, 
  PRIMARY KEY (id)
  );
```
```sql
CREATE TABLE completed (
  id INT NOT NULL,
  name VARCHAR(50) NOT NULL,
  comment VARCHAR(100),
  createDate DATE NOT NULL,
  deadline DATE,
  appointmentDate DATE,
  finishDate DATE,
  PRIMARY KEY (id)
  );
```
2. В файле ```config.yml``` запишите все необходимые данные о базе данных
3. Установить все необзодимые пакеты
4. Запустить cmd/main.go
