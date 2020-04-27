Dash MPD
========

Control MPD from a dash button.

A Dash button sends one of three messages
1. single press
2. double press
3. long press

I think MPD can send messages 
1. when the playlist reaches empty
2. when the player is playing or stops

Actions
-------

### Single Press
* if in playlist mode:
  * play next item
* else:
  * enter playlist mode 
  * play first item
   

### Double Press
* Stop playing everything 
* clear queue
* clear mode

### Long Press
* clear queue
* enter album mode
* Play a random album

### MPD says stop
* if in album mode:
  * clear playlist
* clear mode

### MPD says playlist = 1
* if in album mode:
  * add new random album

