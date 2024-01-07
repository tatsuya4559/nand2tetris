#!/bin/bash
set -eu

go build

test() {
    local -r vmFile="$1"
    local -r testPath="${vmFile%.vm}.tst"

    echo "run ${testPath}"

    ./vmtranslator $vmFile
    ../tools/CPUEmulator.sh $testPath
}

test "../projects/07/StackArithmetic/SimpleAdd/SimpleAdd.vm"
test "../projects/07/StackArithmetic/StackTest/StackTest.vm"
test "../projects/07/MemoryAccess/BasicTest/BasicTest.vm"
test "../projects/07/MemoryAccess/PointerTest/PointerTest.vm"
test "../projects/07/MemoryAccess/StaticTest/StaticTest.vm"

test "../projects/08/ProgramFlow/BasicLoop/BasicLoop.vm"
test "../projects/08/ProgramFlow/FibonacciSeries/FibonacciSeries.vm"
