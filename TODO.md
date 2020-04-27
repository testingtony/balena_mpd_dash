Todo
====

Making it useful for dockerization
- [x] Check whether the client id is required for amazon _yes it is_
- [x] Have it take the amazon connection pems from environment variables
- [x] Have it exit cleanly if the environment variables aren't set

Putting it in docker
- [x] Create a multi-module project
- [x] Add in mpd and dash_mpd as subprojects

Making it better
- [ ] Pick a playlist based on day/time
- [ ] Clear reset when the server sends an unexpected stop message
- [ ] if in album mode and there's on item left in the playlist, add a new album
