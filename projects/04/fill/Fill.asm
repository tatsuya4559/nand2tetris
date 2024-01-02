// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen
// by writing 'black' in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen by writing
// 'white' in every pixel;
// the screen should remain fully clear as long as no key is pressed.

(START)
    // set pointer to end of screen
    @8191
    D=A
    @SCREEN
    D=D+A
    @max_screen_ptr
    M=D
    // start main loop
    @READ_KEY
    0; JMP

(READ_KEY)
    // read keyboard
    @KBD
    D=M
    // if key == 0 then goto CLEAR_SCREEN
    @CLEAR_SCREEN
    D; JEQ
    // else goto FILL_SCREEN
    @FILL_SCREEN
    0; JMP

(CLEAR_SCREEN)
    // ptr := SCREEN
    @SCREEN
    D=A
    @ptr
    M=D
    @CLEAR_LOOP
    0; JMP

(CLEAR_LOOP)
    // clear pixel at ptr
    D=0 // 0 means black (0b_0000_0000_0000_0000
    @ptr
    A=M

    // repeat 32 times for a row
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D

    // ptr += 32
    D=A+1
    @ptr
    M=D
    // if ptr > max_screen_ptr then goto READ_KEY
    // <=> if ptr - max_screen_ptr > 0 then goto READ_KEY
    D=M
    @max_screen_ptr
    D=D-M
    @READ_KEY
    D; JGT
    // else goto FILL_LOOP
    @CLEAR_LOOP
    0; JMP

// fill SCREEN ~ SCREEN + 8K - 1
(FILL_SCREEN)
    // ptr := SCREEN
    @SCREEN
    D=A
    @ptr
    M=D
    @FILL_LOOP
    0; JMP

(FILL_LOOP)
    // fill pixel at ptr
    D=-1 // -1 means black (0b_1111_1111_1111_1111)
    @ptr
    A=M

    // repeat 32 times
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D
    A=A+1
    M=D

    // ptr += 32
    D=A+1
    @ptr
    M=D
    // if ptr > max_screen_ptr then goto READ_KEY
    // <=> if ptr - max_screen_ptr > 0 then goto READ_KEY
    D=M
    @max_screen_ptr
    D=D-M
    @READ_KEY
    D; JGT
    // else goto FILL_LOOP
    @FILL_LOOP
    0; JMP
