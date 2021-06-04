# balance-service
***
### Ручки:
+ Получение информации о балансе:  
Request: **[GET] /{id:[0-9]+}**
  
Response:
<pre>
200
{
  "ID": 1,
  "Balance": 70 
}
</pre>

+ Изменене баланса:  
Request: **[POST] /change-balance**
Body:
<pre>
{
    "id": 1,
    "delta": -10
}
</pre>  

Response:
<pre>
200
{
    "ID": 28,
    "AccountID": 1,
    "CreatedAt": "2021-06-04T10:05:52.7416361Z",
    "Delta": -10,
    "Remaining": 70,
    "Message": "Account [1]: balance changed by [-10.00], [70.00] remaining"
}

403
insuffisient funds on account [1]

400
Validation error(s): 
Key: 'ChangeBalanceRequest.ID' Error:Field validation for 'ID' failed on the 'gt' tag
Key: 'ChangeBalanceRequest.Delta' Error:Field validation for 'Delta' failed on the 'required' tag
</pre>

+ Трансфер:  
  Request: **[POST] /transfer**  
  Body:
<pre>
{
    "id1": 1,
    "id2": 2,
    "delta": 10
}
</pre>  

Response:
<pre>
200
{
    "ID": 24,
    "AccountID": 1,
    "CreatedAt": "2021-06-04T09:13:19.6485027Z",
    "Delta": -10,
    "Remaining": 80,
    "Message": "Transfer from account [1] to account [2]: balance changed by [-10.00], [80.00] remaining"
}

403
insuffisient funds on account [1]

400
Validation error(s): Key: 'TransferRequest.ID1' Error:Field validation for 'ID1' failed on the 'gt' tag
Key: 'TransferRequest.ID2' Error:Field validation for 'ID2' failed on the 'nefield' tag
Key: 'TransferRequest.Delta' Error:Field validation for 'Delta' failed on the 'gt' tag
</pre>

+ Получение истории транзакций:  
  Request: **[GET] /transactions/{id:[0-9]+}?sort=by-time&order=asc&page=2**  
  URL параметры:
  * Параметр **[sort]**  
    Валидные аргументы: *by-time / by-sum*  
    Необязательный параметр, по умолчанию стоит фильтр по времени
  * Параметр **[order]**  
    Валидные аргументы: *asc / desc*  
    Необязательный параметр, по умолчанию стоит фильтр по возрастанию
  * Параметр **[page]**  
    Необязательный параметр, при отрицательных значениях выдает весь список транзакций, 
    при *page=0* равен 1, по умолчанию равен -1

Response:
<pre>
200
[
    {
        "ID": 26,
        "AccountID": 1,
        "CreatedAt": "2021-06-04T09:19:31.616356Z",
        "Delta": -10,
        "Remaining": 70,
        "Message": "Account [1]: balance changed by [-10.00], [70.00] remaining"
    },
    ...
    {
        "ID": 32,
        "AccountID": 1,
        "CreatedAt": "2021-06-04T10:05:52.741636Z",
        "Delta": -10,
        "Remaining": 70,
        "Message": "Account [1]: balance changed by [-10.00], [70.00] remaining"
    }
]

400
Query param [sort] not valid: valid options are [by-sum], [by-time]
</pre>

***

### Переменные конфига:

+ SERVER
    * PORT - порт на котором запускается сервер 
+ DB
    * HOST - хост БД 
    * USER - пользователь БД 
    * PASSWORD - пароль БД 
    * NAME - имя БД 
    * PORT - порт БД 
    * SSL - режим SSL БД
+ SETTINGS
    * PAGINATION_NUM - количество транзакций на странице
    
***

### Запуск:

Для запуска нужно создать PostgreSQL БД с табличками из *balance_tables.sql*. 
В примере ниже БД создается в Docker контейнере с именем pg_balance
+ docker build . -t balance_srv
+ docker run --link pg_balance --rm -p 8081:8081 -d --name balance balance_srv balance-service
+ docker kill balance