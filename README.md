## Solary

A turn based game for bots.


### Rules

In every turn each player can move one field in any direction ("left", "right", "up", "down"), and can use an item from its inventory. Loot is automatically collected if you are the only player on the field.


### Goal

Reach higher score than your opponents


### Items

 - Pogo stick: jump two fields
 - Trap: place a trap to a field, which deals 50 damage on the given field if activated (note: without direction it instantly activates)
 - Laser beam: Deals 25 damage to two fields in the given direction (note: without direction it deals 50 damage to your field)
 - Oil: Heals 20
 - Solar panel: cannot use, produces 1 score per round


### Protocol

Each server/client message should be a single line of json serialized data


#### Player move

JSON object with `direction` and `item` optional keys


##### Examples:

 - `{"direction": "up"}` - move up
 - `{"direction": "right", "item": "laser beam"}` - shoot and move left


#### Server status

JSON object with player and board information
