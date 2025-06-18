Для запуска проекта на linux выполните:
1) скопируйте содержимое файла makefile; создайте файл с названием makefile; вставьте содержимое в makefile вашей директории
2) перейдите в https://github.com/Senkoker/SSO_service/tree/main и выполните указанные в файле readme команды
3) перейдите в директорию PostCommentAuthorization и введите make 

Функционал проекта:
Для более лучшей демонстрации написаны простенькие html страницы с js кодом (путь ./internal/source/sso_html/sources). Также будут показаны запросы в Postman. На данный момент проект переделывается под graphQL c использованием kafka.

Регистрация (смотри html - register.html)
![image](https://github.com/user-attachments/assets/c3aaccd1-81cb-4490-8f00-82d3ef987dcb)

Так как для авторизации и аутентификации используется gRPC_SSO_service, смотрите подробнее https://github.com/Senkoker/SSO_service/tree/main

На прошлом шаге мы уже зарегистрировали пользователя, теперь получаем код ошибки 409

![image](https://github.com/user-attachments/assets/cfed3eb1-f921-4a1e-8c45-96019563a520)

Получаем подтверждение регистрации на почту
![image](https://github.com/user-attachments/assets/ed37d1c7-3b26-408f-a618-e33651d244cf)

Подтверждаем 
![image](https://github.com/user-attachments/assets/3501a8eb-c6e6-4774-af6a-407c3c50bce7)

Логинимся (смотри html - login.html)

![image](https://github.com/user-attachments/assets/80221686-b0fa-4d37-9a36-324255bf2e65)

Через Postman

![image](https://github.com/user-attachments/assets/1543cccb-2908-4342-bd63-9bce3300b8c0)


Заполняем данные пользователя (смотри html - filluser.html). Отправляем данные через formData, так как загружаем изображение.

![image](https://github.com/user-attachments/assets/aaab96a6-d417-46e1-bf7c-2fc9e32ab88f)

fetch запрос через js 

![image](https://github.com/user-attachments/assets/d77a9dbf-f31f-4d14-8708-06d533a89099)

Создаем пост (смотри html - feedExperimental.html). Отправляем данные через formData, так как загружаем изображение.

![image](https://github.com/user-attachments/assets/5202ad0a-95bb-47a7-a1ad-fc23e9328fcc)

![image](https://github.com/user-attachments/assets/b5116ba8-8888-4f86-af19-1bd084f92c0c)

копируем полученный id 
![image](https://github.com/user-attachments/assets/4868a5fb-eea2-4777-853c-5d5ede3fb222)

полученный id вставляем в поиск, получаем 

![image](https://github.com/user-attachments/assets/1f6072f8-0b5e-47bc-a962-376213394eb8)

Данный этап состоит из получения данных пользователя и получения информации о постах(смотрите ниже)

![image](https://github.com/user-attachments/assets/7942d647-042f-4555-9550-bd64b16c58d5)

Можно получить информацию сразу о нескольких постах. Ниже указан поиск по нескольким id 

![image](https://github.com/user-attachments/assets/10653b8b-5f3e-4608-937c-50ecd5e4cdee)

Поиск по нескольким hashtags

![image](https://github.com/user-attachments/assets/81d49bd3-c314-42d8-b8e3-71312be9f04e)

Если посмотрим на параметры, то увидим redis=false, это значит мы напрямую берем данные из Postgres, redis = true - достаем из redis посты. Это работает только для поиска постов по hashtags (поиск в Redis по hastags через ft.Search).Я решил сделать именно так для поиска по hashtags, потому что пользователь ищет популярные посты, а они как раз находятся в кэше. Если при redis = true нет постов, то фронтенд делает запрос redis = false (что поделать придется подождать)

Для поиска постов по ID используется путь: идем в кэш, если в кэше нет, то идем в postgres, достаем из Postgres посты, записываем посты в redis ( прогреваем кэш).   

Давайте посмотрим как прогревается кэш 

Создаем пост 
![image](https://github.com/user-attachments/assets/d5efd305-8061-4238-ba21-b58e1c9288b4)

Пока его нет в Redis (6 постов)

![image](https://github.com/user-attachments/assets/71fac38d-25ee-4c8a-9fc5-c04bc2a6a463)

Делаем get запрос по посту 

![image](https://github.com/user-attachments/assets/90ad86a3-14ce-477b-9338-270048fa1a23)

Смотрим в redis, теперь постов 7

![image](https://github.com/user-attachments/assets/ee715ca3-a01a-4b07-b98b-e45bc3c5b9c5)

Создание комментариев 

Создание комментария родителя (параметр main = true, создание комментариев child через параметр main = false), не получится добавить рандомных пользователей установлен constraint ( фото ниже) 

![image](https://github.com/user-attachments/assets/88b35faf-5185-4b28-b31b-c54bfdb8ff21)

![image](https://github.com/user-attachments/assets/31650bc3-2d54-4a38-97f5-d5694281c3c3)

Создание комментария child

![image](https://github.com/user-attachments/assets/3696390b-7d86-42b2-b53f-0a306eb93e02)

Комментарии не древовидной структуры, то есть главный комментарий и его "дети", которые должны быть упорядочены по времени, так сохранится логика комментариев вопросов и ответов. 

Получение комментариев main

![image](https://github.com/user-attachments/assets/e43d5c6e-b908-4feb-9336-db2e5403068b)

Получение комментариев child

![image](https://github.com/user-attachments/assets/d0c27efb-e7ee-4498-8144-362d9bcddd72)

Проделав данную работу, я понял, что проекты, где используется большое количество связанной информации лучше реализовать на graphQL, что позволит не создавать "миллион" ручек под каждую опцию, к тому же легко создавать вложенные данные (комментарии, посты, информация о пользователе - группы, посты, друзья). Использование последней приводит к новым проблемам: проблема N+1, трудность реализации аутентификации пользователя.
