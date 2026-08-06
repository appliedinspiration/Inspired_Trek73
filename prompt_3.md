Change #1
I am now seeing the damage reports when a ship takes damage, however the timing of these messages is not correct. There is also a problem with the messages indicating that my ship has fired torpedoes...it is being displayed one turn too late.  

Here is an excerpt from a game:
Command: 2
   Fire which tube(s) [all or digits]? 12
Sulu:  Meteor attacking.

Command: 2
   Fire which tube(s) [all or digits]? 34
 <<Potempkin frng torpedo 3>>
 <<Potempkin frng torpedo 4>>
Sulu:  Meteor attacking.

The "potempkin frng torpedo 3" a "4" SHOULD have been displayed immediately after I issued command #2 (firing tubes 1 and 2), but instead it was displayed after I issued command #3 (firing tubes 3 and 4). This indicates that the messages for firing torpedoes are being delayed by one turn, which is not the intended behavior. 

Second, the timing of the messages telling of damage to a ship that was just added such as:
hit 110 on Meteor's shield 1
** probe 9 **
hit 113 on Meteor's shield 1
** probe 9 **

Are not being displayed as they happen. Either they are being delayed by one or more turns, or else they are not being displayed until the game is "over" (after the player has won or lost). This is not the intended behavior either. The damage reports should be displayed immediately after the damage occurs, so that the player can see the results of their actions in real-time.

Please fix these issues.