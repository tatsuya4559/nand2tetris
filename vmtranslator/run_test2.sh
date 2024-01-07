#!/bin/bash
set -eu

go build

test() {
    local -r testPath="$1"
    local -r stem="${testPath%.tst}"
    local -r testDir="${testPath%/*}"

    echo "run ${testPath}"

    ./vmtranslator $testDir
    ../tools/CPUEmulator.sh $testPath || git diff --no-index "${stem}.cmp" "${stem}.out"
}

test "../projects/08/FunctionCalls/FibonacciElement/FibonacciElement.tst"
test "../projects/08/FunctionCalls/NestedCall/NestedCall.tst"
test "../projects/08/FunctionCalls/StaticsTest/StaticsTest.tst"
