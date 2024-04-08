Multiplayer Rock Paper Scissors Betting Game
============================================

How to run
----------
./mrps-game serve --dsn 'host=localhost user=your_username password=your_password dbname=mrps_game sslmode=disable TimeZone=Asia/Shanghai'

Requirements
------------

- [x] Login the system and logout
- [x] Transfer imaginary funds in and out of your account.
- [x] Query for other players
- [x] Challenge other players to a match and associate a choice (rock, paper or scissor) and a bet with it
- [x] View your pending challenges and respond to them with decline or accepting. Accepting requires a choice of rock, paper, scissor
- [x] Upon accepting a challenge, the match is resolved and the betted funds are transferred from the loser to the winner.

Extra tasks:
- [x] Users require registration and this is kept in some kind of database
- [x] All transactions are logged in a database
- [x] Query for personal transactions
- [ ] AI players
- [x] Users are somehow notified if they win or lose
- [ ] Statistics are generated and can be viewed somehow
- [x] Available choices should be configurable (add more than rock paper scissor)

