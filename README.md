# vk-status-change
Установка в статус ВКонтакте данные из Steam
![image](http://image.prntscr.com/image/b6a3982fbde94d0db083211308d7dd5f.png)
# Запуск
Для запуска требуются данные:
* Steam Token
* SteamID64
* VkToken

### Steam Token
https://steamcommunity.com/dev/apikey
### SteamID64
https://steamid.io/
### VkToken
Для получение токена ВКонтакте используйте ссылку:
```
https://oauth.vk.com/token?grant_type=password&client_id=2274003&client_secret=hHbZxrka2uZ6jB1inYsH&username=ЛОГИН&password=ПАРОЛЬ
```
Вместе **логина** и **пароля** вставьте соответствующие данные.

Полученные данные введите в файл **settings.ini** с помощью блокнота.
