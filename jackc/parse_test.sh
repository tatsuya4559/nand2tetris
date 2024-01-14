#!/bin/bash
set -eu

parse() {
    local -r path="$1"

    dune exec jackc -- $path
}

parse ../projects/10/Square/Main.jack
parse ../projects/10/Square/Square.jack
parse ../projects/10/Square/SquareGame.jack
parse ../projects/10/ArrayTest/Main.jack
parse ../projects/10/ExpressionLessSquare/Main.jack
parse ../projects/10/ExpressionLessSquare/Square.jack
parse ../projects/10/ExpressionLessSquare/SquareGame.jack
