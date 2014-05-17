MafiaChat
=========

State “Lobby”, Observing or joining game:
-----------------------------------------

After connecting to server with websocket, server will send following package:

    {
      msgType: “gameInit”,
      data: {
        players: [
          { name: “jack”, role: “villager” },
          { name: “john”, role: “villager” }
        ],
        game: {
          name: “Friday Night Mafia”,
          running: false
        }
      }
    }

If the game is not running, client can join the game by sending a “joinGame” message, with "password" to rejoin the same game and "joinAs" set to "player":

    {
      msgType: “joinGame”,
      data: {
        name: “jane”,
        password: “pass1”,
        joinAs: “player”
      }
    }

If player wants to observe, or the game is already running, joinAs must be “observer”:

    {
      msgType: “joinGame”,
      data: {
        name: “jane”,
        password: “pass1”,
        joinAs: “observer”
      }
    }
    