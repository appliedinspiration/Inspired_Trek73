We are going to add a new damage event to the game, which is a ship detonation due to damage. This damage event mimics the damage caused by a ship that detonates due to a successful self-destruction sequence. However, this detonation can occur due to damage to the shp. The purpose of this new damage event is to allow for a "chain reacton" of damage and ship detonations due to battle damage.  

When a detonation due to damage occurs, damage to surrounding objects is computed the same as if the ship had detonated due to a successful self-destruct. The same messages will be displayed (e.g. "++trakka++ destruct.") and the ship will be listed as destroyed at the end of the game. Damage will be computed and applied to surrounding objects in the same manner as a self-destruct detonation.


When a shp receives damage, whether or not it detonates due to damage is determined by the following procedure:
- If the hit is taken on a shield that is at 100% strength, the ship will not detonate due to damage.
- In the text below, "hit" refers to a hit that penetrates the shields and causes damage to ANY part of the ship.
- If a hit penetrates the shields on a ship whose engines/warp drive are at 100% strength (that is, they have taken no damage), the ship will not detonate due to damage.
- If a hit that penetrates the shields is taken on a ship whose engines/warp drive are at least 75% damaged, the percentage chance of a detonation due to damage is as follows:
-    (damage from the hit taken) / 4
-   For example, if a ship takes 40 points of damage from a hit that penetrates the shields, the chance of detonation due to damage is 40 / 4 = 10%. If the ship takes 90 points of damage from a hit that penetrates the shields, the chance of detonation due to damage is 90 / 4 = 22.5%.