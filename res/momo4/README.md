# Momodora: Reverie Under the Moonlight tracker

The table bellow list all in-game resources and their ID/type in the tracker.

Note that Vitality Framents and Ivory Bugs have two separated resource. The boolean resources (`54` and `34`) are used to light up to icon/image, while the text resource (`54-text` and `34-text`) are used to set the actual value. Therefore, when setting these values to any non-zero value a `POST` should be sent to the boolean resource, and when settign them to zero a `DELETE` should be sent.

| Resource | ID | Type |
| -- | -- | -- |
| Adorned Ring | 1 | boolean |
| Necklace Of Sacrifice | 2 | boolean |
| Bellflower | 4 | boolean |
| Astral Charm | 5 | boolean |
| Edea Pearl | 6 | boolean |
| Dull Pearl | 7 | boolean |
| Red Ring | 8 | boolean |
| Magnet Stone | 9 | boolean |
| Rotten Bellflower | 10 | boolean |
| Faerie Tear | 11 | boolean |
| Impurity Flask | 13 | boolean |
| Passiflora | 14 | boolean |
| Crystal Seed | 15 | boolean |
| Medal Of Equivalence | 16 | boolean |
| Tainted Missive | 17 | boolean |
| Black Sachet | 18 | boolean |
| Ring Of Candor | 21 | boolean |
| Small Coin | 22 | boolean |
| Backman Patch | 23 | boolean |
| Cat Sphere | 24 | boolean |
| Hazel Badge | 25 | boolean |
| Torn Branch | 26 | boolean |
| Monastery Key | 27 | boolean |
| Clarity Shard | 31 | boolean |
| Dirty Shroom | 32 | boolean |
| Violet Sprite | 35 | boolean |
| Soft Tissue | 36 | boolean |
| Garden Key | 37 | boolean |
| Sparse Thread | 38 | boolean |
| Blessing Charm | 39 | boolean |
| Heavy Arrows | 40 | boolean |
| Bloodstained Tissue | 41 | boolean |
| Maple Leaf | 42 | boolean |
| Fresh Spring Leaf | 43 | boolean |
| Pocket Incensory | 44 | boolean |
| Birthstone | 45 | boolean |
| Quick Arrows | 46 | boolean |
| Drilling Arrows | 47 | boolean |
| Sealed Wind | 48 | boolean |
| Cinder Key | 49 | boolean |
| Fragment: Bow Lv.3 | 50 | boolean |
| Fragment: Bow Quick Charge | 51 | boolean |
| Fragment: Dash | 52 | boolean |
| Fragment: Warp | 53 | boolean |
| Vitality Fragment (icon) | 54 | boolean |
| Vitality Fragment (number) | 54-text | string |
| Ivory Bug (icon) | 34 | boolean |
| Ivory Bug (number) | 34-text | string |
