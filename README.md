# Gringoats
## The pokebank backend you always wanted

![stability-wip](https://img.shields.io/badge/stability-wip-red.svg)
----
### What do you mean by pokebank backend ?
Gringoats aims to provide a viable solution to store your pokemons from every game ever created in the same space. Your banks are safe and replicated in your own server and available at any moment.

### Pokemon support

| Pokemon version                             | Compatibility   | Conversion  | Comments                         |
|---------------------------------------------|-----------------|-------------|----------------------------------|
| Red, Blue, Green, Yellow                    | Almost complete | To 2G games | Missing pokemon/player nicknames |
| Silver, Gold, Crystal                       | Almost complete | To 1G games | Missing pokemon/player nicknames |
| Ruby, Sapphire, FireRed, LeafGreen, Emerald | WIP             | No          |                                  |

### How to build and start
```shell
go get ./...
go build
./Gringoats
```
