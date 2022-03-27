# **TestTask** 

## About:
My name is Dmitry and it is my tesk assignment. 

A service is implemented by goLang using socket as a client.

- Each player is able to create a roomsfor a given game.
- Each player is able to join to an existing room if max number of players per room is not reached.
- Maximum 4 players are able to join one room.
- All the actions made by players are visible to each other within the room.

## Requirements:
- Docker
- Docker-compose

# **Run**
1. Write your MySQL environment variables in `docker-compose.yml` and `config/config.yaml`. However all instructions will run without any manipulations:

    - docker-compose.yml

        ```
        ...
        MYSQL_DATABASE: 'game_chat_db'
        MYSQL_USER: 'game_chat'
        MYSQL_PASSWORD: 'password'
        MYSQL_ROOT_PASSWORD: 'password'
        ...
        ```

    - config/config.yaml

        ```
        db:
            host: "db"                  # keep
            port: "3306"                # keep
            database: "game_chat_db"
            user: "game_chat"
            password: "password"
        ```

2. Run docker:

    ```
    $ docker-compose -f docker-compose.yml up
    ```

3. Next you can test API with Postman (Unfortunately, I didn't have enough time to make the front-end for the service and create swagger documenation)

## **API**

* `localhost:8080/` - (`GET`) return list of available players, rooms and games 
* `localhost:8080/login/{number}` - (`POST`) set UserID cookie to authentication, where 'number' - player's id  
* `localhost:8080/logout` - (`POST`) delete UserID cookie (need UserID cookie) 
* `localhost:8080/room` - (`POST`) create room with game. (need UserID cookie and JSON in body 
    ```
    {
        "id": "game_id"
    }
    ``` 
* `localhost:8080/room/uuid` - (`websocket connection`) connect to the room 