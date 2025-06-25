# pug compiler generated assembly
.section __DATA,__data

.section __TEXT,__text,regular,pure_instructions
.globl _main

_main:
    pushq %rbp
    movq %rsp, %rbp
    subq $256, %rsp
    movq $10, %rax
    movq %rax, -8(%rbp)
    # let x = ...
    movq $20, %rax
    movq %rax, -16(%rbp)
    # let y = ...
    movq $5, %rax
    movq %rax, -24(%rbp)
    # let z = ...
    movq -8(%rbp), %rax
    # load variable x
    pushq %rax
    movq -16(%rbp), %rax
    # load variable y
    movq %rax, %rbx
    popq %rax
    addq %rbx, %rax
    pushq %rax
    movq -24(%rbp), %rax
    # load variable z
    movq %rax, %rbx
    popq %rax
    imulq %rbx, %rax
    pushq %rax
    movq -8(%rbp), %rax
    # load variable x
    movq %rax, %rbx
    popq %rax
    subq %rbx, %rax
    movq %rax, -32(%rbp)
    # let result = ...
    movq -32(%rbp), %rax
    # load variable result
    movq %rbp, %rsp
    popq %rbp
    ret
    movq $0, %rax
    movq %rbp, %rsp
    popq %rbp
    ret
