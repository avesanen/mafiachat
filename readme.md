MafiaChat
=========

Quick ideas after DevX event
----------------------------


Chatroom based on Mafia partygame.

Root view:
- Games running
- Games not started yet
- Create new game

Game views:
    Villager:
    - Chat with only public text showing
    - Vote and change vote until timer runs out, or majority votes same person.
    Mafioso:
    - Chat with public text showing
    - Private channel (same channel, only hidden from others, different color)
    - During night, vote to kill someone until timer runs out or majority (all?) votes match.
    Dead:
    - Chat only with other dead players
    - See all private messages
  Later:
    Doctor:
    - Select player to protect, change protected player
    - Chat on public
    - Doctor private messages
    Police:
    - Select player to know if he is mafioso once night ends.
    Crazy:
    - Select player to kill once night ends.


*Server: Day 1, daytime. Night starts in 30 minutes.
<player1> I'm not the bad guy!
<player2> sure you are!
<player3> /vote player1
<player2> /vote player1
*Server player1 has majority vote, 30 seconds to live unless votes change
<player1> REALLY! I'm not! 3 has to be!
<player1> /vote player3
<player2> fine, i believe you
<player2> /vote player3
*Server: player1 lost majority vote
*Server: player3 has majority vote, 30 seconds to live unless votes change
<player1> bye bye
*Server: player3 has been eliminated.
*Server: Nighttime starts, players muted, mafioso can chat and vote.