@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_0
D;JNE
@SP
A=M-1
M=-1
(END_SET_TRUE_0)
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_1
D;JNE
@SP
A=M-1
M=-1
(END_SET_TRUE_1)
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_2
D;JNE
@SP
A=M-1
M=-1
(END_SET_TRUE_2)
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_3
D;JGE
@SP
A=M-1
M=-1
(END_SET_TRUE_3)
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_4
D;JGE
@SP
A=M-1
M=-1
(END_SET_TRUE_4)
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_5
D;JGE
@SP
A=M-1
M=-1
(END_SET_TRUE_5)
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_6
D;JLE
@SP
A=M-1
M=-1
(END_SET_TRUE_6)
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_7
D;JLE
@SP
A=M-1
M=-1
(END_SET_TRUE_7)
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@SP
A=M-1
M=0
@END_SET_TRUE_8
D;JLE
@SP
A=M-1
M=-1
(END_SET_TRUE_8)
@57
D=A
@SP
A=M
M=D
@SP
M=M+1
@31
D=A
@SP
A=M
M=D
@SP
M=M+1
@53
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=D+M
@112
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=M-D
@SP
A=M-1
M=-M
@SP
AM=M-1
D=M
A=A-1
M=D&M
@82
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=D|M
@SP
A=M-1
M=!M
