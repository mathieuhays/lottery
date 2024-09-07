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

- [x] Improve styling
- [ ] Add a "multi-thread" mode if possible
- [ ] add option to restart the brute force process
- [ ] add option to reset the simulation and give opportunity to input new values
- [ ] switch interval into option selector with predefined options
- [ ] interval selector, add option to change interval during brute force process?
- [ ] add notification when winning if technically possible
