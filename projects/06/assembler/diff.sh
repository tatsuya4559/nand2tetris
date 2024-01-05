#!/bin/bash
set -eu

get_diff() {
    local -r asm="$1"
    local -r stem="${1%.asm}"
    local -r hack="${stem}.hack"
    local -r got="${stem}-got.hack"
    local -r want="${stem}-want.hack"

    rm -f $hack

    ./assembler $asm
    mv $hack $got

    ../../../tools/Assembler.sh $asm
    mv .hack $want

    git diff --no-index $want $got
    # vimdiff $want $got
    echo OK
}

get_diff $1
