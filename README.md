[![CircleCI](https://circleci.com/gh/h4ck3rk3y/go-8.svg?style=svg)](https://circleci.com/gh/h4ck3rk3y/go-8)

# go-8 
A chip-8 emulator written in Go. I am trying to learn go-lang for fun and I have always been fascinated with emulators, so I decided to write one in Go. I am using the following [manual](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#1nnn) as reference.

You'll need the drivers required by ebiten and you can find the relevant installation instructions [here](https://github.com/hajimehoshi/ebiten/wiki/Linux).

## Test Instructions

```bash
- git clone http://github.com/h4ck3rk3y/go-8
- cd go-8
- go get -v -t -d ./...
- go test --cover
```

## Build Instructions

You can build the program by

```bash
- git clone http://github.com/h4ck3rk3y/go-8
- cd go-8
- go get -v -t -d ./...
- go build main.go cpu.go
```

## Run Instructions

Build the code and then

```bash
- ./main
```
## Key Configuration

The original chip-8 consisted of a hexa decimal gamepad. I use the following mappings.

- Your Key Board --> Chip 8
- 1 --> 1
- 2 --> 2
- 3 --> 3
- 4 --> C
- Q --> 4
- W --> 5
- E --> 6
- R --> D
- A --> 7
- S --> 8
- D --> 9
- F --> E
- Z --> A
- X --> 0
- C --> B
- V --> F

## To Do

- Make roms passable as command line arguments
- Key board mapping in a configuration file
- Configurable colors
- Better unit tests for main.go. cpu.go has 98.8% coverage but overall the coverage drops significantly
