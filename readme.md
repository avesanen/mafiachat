MafiaChat
=========

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
