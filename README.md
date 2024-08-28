# Anthophila
## У стадії розробки

[Logo](https://github.com/rifatismailov/Anthophila/blob/master/Anthophila.gif)
[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://github.com/rifatismailov/Anthophila)

Anthophila – програма для моніторингу та збору інформації.

- Для запуску виконуємо команду ./Anthophila -file_server="localhost:9090" -manager_server="localhost:8080" -log_server="localhost:7070" -directories="?" -extensions=".doc,.docx,.xls,.xlsx,.ppt,.pptx" -hour=12 -minute=45 -key="a very very very very secret key" -log_file_status=true -log_manager_status=true

- і чекаємо

## Мета розробки:

- автоматизований пошук файлів на комп'ютері та передача їх на сервер
- відправка файлів на віддалений сервер для пошуку критичної інформації


Програма виконує пошук докуменів формату DOC, DOCX, Excel, PowerPoint та інші документи які вам потрібні для аналізу, а також надає можливість віддаленого доступу адміністратора до комп'ютера.


## Основні функції:

- Пошук файлів та передача їх на сервер.
- Зберігання Hash суму файлів для унеможливлення повторного відправлення на сервер тільки якщо файл був змінений. 
- Передача файлів та повідомлення у зашифрованому вигляді (стосовно повідомлень ще не реалізована)
- Надання віддаленого доступу адміністратору до комп'ютера через CMD
- Можливість копіювання файлів з комп'ютера на сервер за командою адміністратора.(ще не реалізована)


## Інсталяція або запуску

Для запуску Anthophila потрібно ваші не криві руки та люба операційна система.

Інсталяція не потрібна. Запускається з компільованого файлу або ви самі можете завантажити вихідний код та скопілювати в себе.
Команда для компіляції
```sh
./Anthophila -file_server="localhost:9090" -manager_server="localhost:8080" -log_server="localhost:7070" -directories="?" -extensions=".doc,.docx,.xls,.xlsx,.ppt,.pptx" -hour=12 -minute=45 -key="a very very very very secret key" -log_file_status=true -log_manager_status=true
```

На що вам треба звернути увагу під час запуску программи
  
```sh
1. Правильно вказати все три сервера (Файл сервер, Лог сервер та сервер керування)
2. Правильні діректорії які вам треба сканувати. Приклад :
   [Вона може бути виглядати так -directories= "/Users/username/Desktop/,/Users/username/Documents/, /Users/username/Downloads/"
3. Типи файлів у такому вигляді -extensions=".doc,.docx,.xls,.xlsx,.ppt,.pptx"
4. Час коли программа починає сканування -hour=12 -minute=45
5. Ключ шифрування -key="a very very very very secret key"
6. Можливість віддаленого логування тих помилок які вам потрібні. Допомогає коли програма запушена на машині і дає можливість збору логів помилок для подальшого їх виришення. 

```

## Сумісність:

- Програма працює на всіх операційних системах в залежності від бінарного файлу який буде скопільований для кожної оперційної системи.

## Інтерфейс користувача:

- Програма працюватиме у фоновому режимі без графічного інтерфейсу.

Програма написана на Go



## Ліцензія

**Безкоштовне програмне забезпечення!**


Insert gif or link to demo
![Logo](https://github.com/rifatismailov/Anthophila/blob/master/Anthophila.gif)