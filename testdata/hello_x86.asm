; === [ text section ] =========================================================

section ".text"
global _start

_start:
	; write(1, "Hello world\n", 12)
	mov	eax, 4		; sys_write
	mov	ebx, 1		;    fd: stdout
	mov	ecx, hello	;    buf: str
	mov	edx, 12		;    len: len(str)
	int	0x80		; syscall

	; exit(0)
	mov	eax, 1		; sys_exit
	mov	ebx, 0		;    status = 0
	int	0x80		; syscall

; === [ rdata section ] ========================================================

section ".rdata"

hello:	db "Hello world", '!', 10
