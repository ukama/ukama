OUTPUT_FORMAT("elf32-littlearm", "elf32-littlearm", "elf32-littlearm")
OUTPUT_ARCH(arm)
ENTRY(reset)
SECTIONS
{
        . = ALIGN(4);
	.text : {
		_stext = .;
		*(EXCLUDE_FILE (*mon_*.o) .text)
		*(EXCLUDE_FILE (*mon_*.o) .rodata)                 /* read-only data (constants) */
		*(EXCLUDE_FILE (*mon_*.o) .rodata*)
		. = ALIGN(4);
		*(EXCLUDE_FILE (*mon_*.o) .glue_7)
		. = ALIGN(4);
		*(EXCLUDE_FILE (*mon_*.o) .eh_frame)
		. = ALIGN(4);
 		_etext = . ;
	}

	. = ALIGN(4);
	.dummy : {
		_edummy = .;
	}

	.data : AT (LOADADDR(.dummy)) {
		_sdata = .;
		*(EXCLUDE_FILE (*mon_*.o) .vectors)
		*(EXCLUDE_FILE (*mon_*.o) .data)
		_edata = .;
	}

	/* collect all uninitialized .bss sections */
	.bss (NOLOAD) : {
		. = ALIGN(4);
		_sbss = .;
		*(EXCLUDE_FILE (*mon_*.o) .bss)
		_ebss = .;
	}

	. = ALIGN(4);
	.mon_dummy : {
		_emon_dummy = .;
	}

	. = ADDR(.text) + 0x8000;
	.mon_text : AT (_emon_dummy) {
		_smon_text = .;
		*mon_vectors.o		/* First input */
		*mon_init.o
		*mon_switch.o
		_emon_text = .;
	}
}
_romsize = _edata - _stext;
_sramsize = _emon_text - _stext;
end = .;  /* define a global symbol marking the end of application */
