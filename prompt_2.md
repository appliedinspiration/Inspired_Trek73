Change 1:
Add the following functionality: 
whenever a ship takes damage (and this is specific to ships only), a message should be output of the form:
hit <damage amount> on <ship>'s shield <shield number>
So if the Lexington sustains a 45 point hit on shield #2 the message should read:
"hit 45 on Lexington's shield 2"

Change 2:
At program start up, after the program asks for the captain's name and also the captain's sex, if no command line argument is specified for the number of enemy vessels, ask the user the following question: "I'm expecting [1-9] enemy vessels: " and let them enter a value between 1 and 9. If the user enters a value outside of this range, prompt them again until they enter a valid number. 

Change 3:
The position display (option #13) is displaying an ascii image that does nto look correct. Please fix the ascii image so that it displays correctly in the terminal. There does not appear to be enough "vertical lines" in the displayi so that position display always looks like it has been "flattened" or "squished". Please adjust the ascii image to ensure that it maintains its intended proportions and looks correct in the terminal.

Change 4:
The "Commands" menu option #32 lists the avaialable commands in one long vertical list. Please change this to display the listing in two columns rather than just one. Also, if at the commands menu the user enters a blank entry rather than a number, assume that they meant to enter 32 and redisplay the commands list (rather than displaying an error message).