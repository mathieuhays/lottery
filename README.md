# Lottery Simulator

uses a brute-force approach to demonstrate the odds of winning at the lottery

## Build

```shell
go build -o lottery
```

## Usage

```shell
./lottery
```

You will be prompted with field to specify the odds, the interval to draw tickets at and the cost per ticket.

Example settings:
- Chances: 25 000 000 (as in 1/25 000 000 changes to win)
- Interval: 100 (in ms). draw 10 tickets a second
- Cost: 0.25 (no units)

Given those parameters you could expect the script to draw 10 tickets a second for 30 days to match the odds (no guaranteed win within those odds) with an associated cost of 6,250,000.

## TO DO

- [ ] Mode: Guaranteed odds
- [ ] Mode: buy mode
- [ ] Improve styling
- [ ] Add a multi-thread mode if possible
- [ ] Brute-force mode: add option to restart the brute force process
- [ ] Brute-force mode: switch interval into option selector with predefined options
- [ ] Brute-force mode: interval selector, add option to change interval during brute force process?

### Guaranteed odds mode

The odds represent the maximum amount of tickets "purchasable". Each ticket can only be drawn once.

Potential memory consumption involved with storing a hashmap of unique numbers would be around 100mb for 25,000,000 numbers assuming int 32 is used. Up to 2Gb for 500,000,000 numbers. With a max size of around 9Gb for the maximum value of int 32.

Could a tree like structure be used to store drawn numbers in ranges to reduce memory usage? How to draw unique numbers as more of them are drawn and are no longer valid for new tickets. Just calling `rand.Intn` might no longer be feasible without affecting performance at the end of the range.

It could be an option to hide the generated ticket number in that mode and instead reduce the range internally as tickets are drawn so we don't need to keep track of the number used but what range is remaining instead.

### Buy mode

- Buy based on budget, meaning buy as many tickets given a monetary amount
- Buy x amount of tickets
- Add win animation with particles or something?