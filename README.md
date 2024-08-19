# Docker

## build
docker build -t app .

## run
docker run --mount type=bind,source=.,target=/usr/src/app -it app

## other
docker image ls


# GO

# PROTOCOL
## FIELD VALUES
WIDTH = 25
HEIGHT = 25
MINE_VALUE = 10
CLEAN_VALUE = 11

## FIRST MSG
1) Gets empty message

### GAME_STATUSES:
const GAME_WON = 1
const GAME_LOST = 2
const GAME_IN_GAME = 3

2) Returns array, [Width, gamestatus, pointsLength, [player, points...]...field]

