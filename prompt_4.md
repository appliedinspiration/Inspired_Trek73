Item #1
On menu command 12 (the position report), the column headings are not lining up properly with the data columns. In fact, in one case the columns of the 1st row of the report (headings) is not lining up with the second heading row (the columns are not aligned). Please realign these columns.

Item #2
I just played a game where I fired several torpedoes at an enemy ship. The torpedoes hit the target, and the damage report looked like this:
hit 40 on Millennium Pelican's shield 2
hit on torpedo 4
:: torp 4 ::
hit 30 on Millennium Pelican's shield 1
hit on torpedo 3
:: torp 3 ::
hit 42 on Millennium Pelican's shield 2
:: torp 4 ::
hit 33 on Millennium Pelican's shield 1
:: torp 3 ::
hit 42 on Millennium Pelican's shield 2
:: torp 4 ::
hit 37 on Millennium Pelican's shield 1
:: torp 3 ::
hit 42 on Millennium Pelican's shield 2
:: torp 4 ::
hit 39 on Millennium Pelican's shield 2
:: torp 3 ::
hit 40 on Millennium Pelican's shield 2
:: torp 4 ::
hit 41 on Millennium Pelican's shield 2
:: torp 3 ::
hit 38 on Millennium Pelican's shield 3
:: torp 4 ::
hit 42 on Millennium Pelican's shield 2
:: torp 3 ::
hit 35 on Millennium Pelican's shield 3
:: torp 4 ::
hit 42 on Millennium Pelican's shield 2
:: torp 3 ::
hit 32 on Millennium Pelican's shield 3
:: torp 4 ::
hit 41 on Millennium Pelican's shield 2
:: torp 3 ::
hit 28 on Millennium Pelican's shield 3
:: torp 4 ::
hit 39 on Millennium Pelican's shield 3
:: torp 3 ::
hit 23 on Millennium Pelican's shield 3
:: torp 4 ::
hit 36 on Millennium Pelican's shield 3
:: torp 3 ::
hit 17 on Millennium Pelican's shield 3
:: torp 4 ::
hit 33 on Millennium Pelican's shield 3
:: torp 3 ::
hit 8 on Millennium Pelican's shield 3
:: torp 4 ::
hit 29 on Millennium Pelican's shield 3
:: torp 3 ::
:: torp 4 ::
hit 25 on Millennium Pelican's shield 3
:: torp 3 ::
:: torp 4 ::
hit 20 on Millennium Pelican's shield 3
:: torp 3 ::
:: torp 4 ::
hit 13 on Millennium Pelican's shield 3

I can easily see how one or two torpedoes can hit/damage multiple shields, but it looks to me like those two torpedoes are hitting the ship/causing damage way too many time for just two torpedo hits. It seems like the damage report is not accurately reflecting the number of torpedo hits and the corresponding damage to the shields. I would recommend reviewing the damage calculation logic to ensure that each torpedo hit is only counted once per shield, and that the total damage reported aligns with the actual number of torpedo hits.

Please correct these issues or report that they are working correctly if there is nothing wrong.