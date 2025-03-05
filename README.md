# Распределённый вычислитель арифметических выражений
## О проекте
Вычислитель состоит из двух элементов: оркестратор и агент

### Оркестратор
Оркестратор принимает на вход арифметическое выражение и переводит её в набор последовательных задач и обеспечивает порядок их выполнения. Это производится с помощью дерева(ast). У выражения есть своя структура Expression , в которую включены поля для хранения последовательных задач и узлов дерева. Также оркестратор хранит все выражения , которые вы ему отправляете.
У оркестратора есть 5 endpoint-ов:

1) curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": <строка с выражение>
}'    
С помощью данного эндпоинта вы отправляете выражение оркестратору , которое он обрабатывает и сохраняет себе

2) curl --location 'localhost/api/v1/expressions'    
Этот эндпоинт можно использовать для получения всего списка выражений
3) curl --location 'localhost/api/v1/expressions/:id'   
Этот можно использовать для получения конкретного выражения по id , который вы добавили ранее
4) curl --location 'localhost/internal/task' ("GET")       
Этот используется для получения задачи для выполенния
5) curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
"id": 1,
"result": 2.5
}'      
Этот отправляет решенную задачу

Примеры использования эндпоинтов будут дальше

### Агент
Агент запускает несколько горутин , которые задаются переменной среды COMPUTING_POWER, каждая горутина выступает в роли независемого вычислителя.
Каждую секунду агент отправляет запросы оркестратору с помощью ручки GET internal/task. Если он получает задачу , то вычисляет её и отправляет обратно оркестратору с помощью POST internal/task. Если оркестратор отвечает, что в данный момент нет задачи для выполнения , то агент делает паузу на 15 секунд.

## Схема работы

![schema](C:\Users\Владимир\GolandProjects\scheme.png)


Также вот так разбивается выражение в дерево в оркестраторе:

![tree](C:\Users\Владимир\GolandProjects\tree.png)

Схема иллюстрирует то как работает агент и оркестратор между собой , то что я описывал раньше

## Технологии и библиотеки
Вычислитель написан на языке **Go** и использует следующие библиотеки и инструменты:

#### Язык программирования:
- **Go** (версия 1.23.2)


## Структура проекта 

- **`cmd/`**
  - **`agent/`**
    - **`main.go`**
  - **`orchestrator/`**
    - **`main.go`**
- **`internal/`**
  - **`agent/`**
    - **`agent.go`**
  - **`orchestrator/`**
    - **`parser/`**
      - **`parser.go`**
      - **`ast.go`**
      - **`errors.go`**
      - **`node.go`**
    - **`server/`**
      - **`server.go`**
      - **`taskposthandler.go`**
      - **`taskgethandler.go`**
      - **`calchandler.go`**
      - **`expressionshandler.go`**
      - **`expressionhandler.go`**
    - **`test/`**
      - **`calchandler_test.go`**
      - **`taskgethandler_test.go`**
      - **`taskposthandler_test.go`**
  - **`models/`**
    - **`expression.go`**
    - **`expressions.go`**
    - **`badresponse.go`**
    - **`goodresponse.go`**
    - **`id.go`**
    - **`request.go`**
    - **`taskget.go`**
    - **`taskpost.go`** 

- **`pkg/`**
  - **`calculation/`**
    - **`calculation.go`**
    - **`errors.go`**
    - **`calculation_test.go`**
  - **`logger/`**  
    - **`logger.go`**

## Quick start

Для начала нужно склонировать репозиторий командой 
```
git clone git@github.com:VladimirGladky/FinalTaskFirstSprint.git
```

После этого вам нужно перейти в папку с проектом 
```
cd FinalTaskFirstSprint
```

Теперь вы можете запустить проект , но для этогт нужно чтобы был установлен Go версии 1.23.2
Ссылка для скачивания: [Go Download](https://go.dev/doc/install)

Сначала запускается оркестратор , затем запускается агент

```
go run ./cmd/orchestrator/main.go
```

потом открываете другой терминал и запускаете агента

```
go run ./cmd/agent/main.go
```

Для прекращения работы агента или оркестратора можете нажать сочетание клавиш Ctrl+C

## Примеры использования со всеми возможными сценариями

После запуска проекта вы можете отправлять cURL-запросы к сервису:

Так как в терминале Windows не обрабатываются cURL запросы я использовал git bash.

Нужно отметить , что мой веб-сервис использует порт 4040(надеюсь он у вас не занят 😊)

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/calculate' 
```

cURL команда с ответом сервиса 201:

```bash
 curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2"
}'
```
Ответ:

```bash
{"id":"f67287f7-f29f-4196-b3be-7abff3bec739"}
```

cURL команда с ответом сервиса 400:
```bash
curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2
}'
```
Ответ:
```bash
{"error":"Bad request"}
```

cURL команда с ответом сервиса 405:
```bash
curl --request GET \ --url '127.0.0.1:4040/api/v1/calculate' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2"
}'
```
Ответ:
```bash
{"error":"You can use only POST method"}
```

cURL команда с ответом сервиса 422:
```bash
curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2)"
}'
```
Ответ:
```bash
{"error":"Expression is not valid"}
```

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/expressions' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions'
```
Ответ:
```bash
{
  "expressions":
    [
      {
        "id":"56e8677e-a058-485d-bc2c-342af7130c4c",
        "status":"done",
        "result":6
      },
      {
        "id":"7ddf6a72-d2a3-4bcc-947c-998fd9eac383",
        "status":"done",
        "result":10
      }
    ]
}
```

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/expressions/{id}' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions/56e8677e-a058-485d-bc2c-342af7130c4c'
```
Ответ:
```bash
  {
    "id":"56e8677e-a058-485d-bc2c-342af7130c4c",
    "status":"done",
    "result":6}
```

cURL команда с ответом сервиса 404:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions/56e8677e-a058-485d-bc2c-342af7130c4'
```

Ответ:
```bash
{"error":"Expression not found"}
```

### Запросы по endpoint-у:

метод GET

```bash
'127.0.0.1:4040/internal/task' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/internal/task'
```

Ответ:
```bash
{
    "id":"5",
    "arg1":4,
    "arg2":2,
    "operation":"*",
    "operation_time":300
}
```

curl команда с ответом сервиса 404:

```bash
curl --location '127.0.0.1:4040/internal/task'
```

Ответ:
```bash
{"error":"No task available"}
```

Мой телеграм https://t.me/smoothhhhhhh





