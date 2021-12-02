# Gambit

Chess board in your terminal.

<br/>
<p align="center">
  <img width="90%" src="./chess.gif?raw=true" alt="Terminal chess" />
</p>
<br/>

### Warning

`gambit` does not have many features at the moment.
I plan on adding a chess engine, mouse support, timers, networked play, game replays, etc...

### Move

Suppose you want to open as white with Pawn `E4`.
You will first select the square with the piece you want to move by typing `E2`,
then type the square of the desired end position `E4`.

Just like a real chess board, `gambit` currently supports illegal moves.
_This will likely change in the future_.

### Players

`gambit` supports local and networked play. You can play a local game by
running `gambit` and moving the pieces. You can flip the board by pressing
<kbd>ctrl+f</kbd> to allow the second player to go.

For networked play (not available yet), both players can run `gambit unique-room-id`, this will connect
both players to a shared room in which both can take turns making moves.
