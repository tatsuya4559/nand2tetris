// R2 := R0 * R1
(START)
    // mul := 0
    @0
    D=A
    @mul
    M=D
    // count := 0
    @count
    M=D
    // goto LOOP
    @LOOP
    0; JMP

// (mul += R0) * R1 times
(LOOP)
    // D := R1 - count
    @count
    D=M
    @R1
    D=M-D
    // if (R1 - count) <= 0 then goto STORE_RESULT
    @STORE_RESULT
    D; JLE
    // load R0 to D
    @R0
    D=M
    // mul += R0
    @mul
    M=D+M
    // count++
    @count
    M=M+1
    // goto LOOP
    @LOOP
    0; JMP

(STORE_RESULT)
    @mul
    D=M
    @R2
    M=D
    @END
    0; JMP

(END)
    @END
    0; JMP
