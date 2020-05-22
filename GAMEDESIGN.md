# Project Millwheat

In Millwheat you run a village in the 13th century  with the goal of supplying the king with soldiers.
Soldiers can be swordsmen, bowmen or calvary. The type indicates which items they need, next to their need
for one item of wood and one item of drink.
To supply all these items you need to grow your village.

Millwheat is a persistent browser based game which has rounds of battles that are run transparently in the background
but which give the game its competitive loop. After a certain amount of time a soft reset will be done.
Perhaps these 'seasons' can switch the battleground.

Millwheat is inspired on Travian (town sim), Knights & Merchants (themes) and Heroes & Generals (war system).

## Game loop

- Grow your village (logging, mining, farming, blacksmithing)
- Recruit and equip soldiers (with tools and food)
- Supply soldiers and resources for the war effort

You grow your village by constructing buildings and leveling them up. Buildings consume and produce items.
You need certain items to recruit soldiers, namely food and drinks (bread, meat, wine) and armour and weapons
(leather armor, iron chain mail, iron plate armour, longbow, sword, lance).

Every x days the king will organize battles to conquer or reclaim land, you need to help the war effort
by supplying resources but more importantly trained troops.

After every x months a season is completed and every village will be wiped before users can start playing again.

There will be many highscore lists to show user's achievements i.e. town size, town production rate, 
soldiers recruited, battles won etc.

## Buildings, items and professions

These lists are highly inspired by Knights and Merchants.

- Farm (wheat)
- Windmill (wheat -> flour)
- Bakery (flour -> bread)
- Forestry (logs)
- Saw mill (logs -> planks) 
- Iron mine (iron)
- Coal mine (coal)
- Iron smithy (iron + coal -> iron bars)
- Weapon smithy (iron bars, planks -> weaponry)
- Armour smithy (leather, iron bars -> armoury)
- Vineyard (grapes -> wine)
- Pig farm (wheat -> pigs)
- Butchery (pigs -> meat & hide)
- Tannery (hide -> leather)
- Quarry (stone blocks)
- Stables (wheat -> horses)
- Warehouse
- Barracks

Some other possible buildings:

- Charcoal kiln (logs -> (char)coal)
- Fishery (fish)
- Brewery (wheat -> beer)
- Hunter's lodge (meat & hide)
- Winery (grapes -> wine) (changes vineyard to only grow grapes)
- Siege workshop (planks, metal -> ballista, catapult)
- Market (required to trade with other towns)

These are probably not used within the game, but they're here for informational purposes

- Farmer (farm, vineyard)
- Woodcutter (forestry)
- Stonemason (quarry)
- Carpenter (saw mill)
- Miner (iron mine, coal mine)
- Blacksmith (iron smithy, weapon smithy, armour smithy)
- Baker (bakery, windmill)
- Animal breeder (pig farm, stables)
- Butcher (butchery, tannery)


## Simulation example / construction material examples

https://play.golang.org/p/FFr2VJsvLuV

Construction costs for **Farm**

Stone column is calculated by an exponential function while the planks have been adjusted by hand for gameplay value.

|     | Stone | Planks |
|-----|-------|--------|
| 1:  | 1     | 3      |
| 2:  | 2     | 6      |
| 3:  | 3     | 15     |
| 4:  | 5     | 50     |
| 5:  | 7     | 75     |
| 6:  | 12    | 125    |
| 7:  | 19    | 250    |
| 8:  | 30    | 375    |
| 9:  | 48    | 550    |
| 10: | 77    | 750    |
