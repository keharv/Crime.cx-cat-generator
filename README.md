![release workflow](https://github.com/keharv/Crime.cx-cat-generator/actions/workflows/go.yml/badge.svg)
# [Crime.cx](https://Crime.cx) Cat Generator
This web crawler scrapes images of cats from the internet.

These images are then viewable at https://crime.cx and https://crime.cx/images

-Feline fine? This web crawler will make you feel great!

-In the mood for some cat-astrophic fun? This web crawler is perfect!

-Looking for a purr-fectly delightful web crawler? You've found it!


The goal is to create the largest repository of cat images in the world.

## Contributors

This web crawler was lovingly created by the following cat lovers:

-@keharv
-@TheMysteriousMouse
-@immortal-beast
-[@CowSayMoe](https://github.com/cowsaymoe)


## Built With
GoLang
    

[Tinygrad](http://github.com/geohot/tinygrad) with YOLOv2 - Used to determine if there is a cat within the scraped images

## Usage
```bash
./scraper.exe -output="scraped"
```

## Build/Install
```bash
go build -v ./...
go build
```
In the future update repo so that users can just do ```go get https://github.com/keharv/Crime.cx-cat-generator```

## Road Map

- [X] Create multi-threaded bot that scrapes images then visits all links on webpage which stores visited websites in a database
- [X] Store md5 hashes of images in database to ensure no duplicate images are downloaded, even if they are deleted
- [ ] Create bot that removes images that do not contain pictures of cats
- [ ] Obtain additional storage and hosting to hold the future terabytes of cat images
- [ ] Create a better front-end webpage for https://crime.cx/

## Interested in contributing?
If you are interested in contribute, make a pull request! Additionally, you may contact me @ admin@crime.cx !
