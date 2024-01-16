#!/bin/bash
set -eu

make

test() {
    local -r dir="$1"

    cp ../tools/OS/*.vm $dir
    ./jackc $dir
    echo $dir compiled
}

test ../projects/11/Seven/
test ../projects/11/ConvertToBin/
test ../projects/11/Square/
test ../projects/11/Average/
test ../projects/11/Pong/
test ../projects/11/ComplexArrays/
