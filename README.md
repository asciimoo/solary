## Solary

A turn based game for bots.


### Rules

In every turn each player can move in any direction (`left`, `right`, `up`, `down`), and can use an item from its inventory. Loot is automatically collected if you are the only player on the field.

Rocks are not passable however laser beam can go through it.


### Goal

Reach higher score than your opponents after the 500th round.


### Items

 - `pogo stick`: jump two fields
 - `trap`: deals 50 damage to each player on the given field when activated (note: without a direction it instantly activates)
 - `laser beam`: deals 25 damage to two fields in the given direction (note: without direction it deals 50 damage to your field)
 - `oil`: Heals 20
 - `solar panel`: cannot be used, produces 1 score per round


### Protocol

Each server/client message should be a single line of json serialized data.
Default server port is TCP 6666.


#### Player move

JSON object with `Direction` and `Item` optional keys


##### Movement examples

 - `{"Direction": "up"}` - move up
 - `{"Direction": "right", "Item": "laser beam"}` - shoot and move left
 - `{"Item": "oil"}` - stay and heal


#### Server status

JSON object with player and board information
