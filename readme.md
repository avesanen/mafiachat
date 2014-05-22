**MafiaChat**
=========

A cunning game of social engineering and battle of wits.

----

**State "Lobby"**
----

1. **login**
---

After connecting to websocket, client will send login request:

    {
        "msgType":"login",
        "data":{
            "name":"john",
            "password":"passwrd1"
        }
    }
    
If login is successful, server will respond with gameInfo. If client has received gameInfo, it can proceed to chat view.

    {
        "msgType":"gameInfo",
        "data":{
            "id":"deadc0ffeebabe",
            "state": "lobby",
            "players": [
                {
                    "id" : "randomUUID1",
                    "name": "john",
                    "state": "new"
                },
                {
                    "id" : "randomUUID2",
                    "name": "jahn",
                    "state": "new"
                }
            ]
            
        }
    }

If the login is NOT successfull, server will respond with error. In that case, client will display the error and wait for another login attempt.

    {
        "msgType": "error",
        "data": {
            "message": "wrong password for player john"
        }
    }
    
----

2. **Chat**
----

Whenever the game or player states change, gameInit message will be sent again. For example, when player dies, joins the game or disconnects.

**Other messages include:**

*chatMessage:*

    {
        "msgType": "chatMessage",
        "data": {
            "message": "Hello comrades!",
            "faction": "mafia"
        }
    }
    
*serverMessage*

    {
        "msgType": "serverMessage",
        "data": {
            "message": "Player eric has disconnected."
        }
    }
    
*errorMessage*

    {
        "msgType": "errorMessage",
        "data": {
            "message": "Unexpected error has occurred... 0x9215fab232"
        }
    }
    
    