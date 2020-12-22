OUTPUT_FORMAT("elf32-littlearm", "elf32-littlearm", "elf32-littlearm")
OUTPUT_ARCH(arm)
ENTRY(reset)
SECTIONS
{
        . = ALIGN(4);
	.text : { 
		_stext = .;
		*(.text)
		*(.rodata)                 /* read-only data (constants) */
		*(.rodata*)
		. = ALIGN(4);
		*(.glue_7)
		. = ALIGN(4);		
		*(.eh_frame)
		. = ALIGN(4);
 		_etext = . ;
	}

	/* collect all initialized .data sections */
	/* .data : AT ( ADDR (.text) + SIZEOF (.text) SIZEOF (.ARM.*) {  */

	. = ALIGN(4);
	.dummy : {
		_edummy = .;
	}

	.data : AT (LOADADDR(.dummy)) {
		_sdata = .;
		*(.vectors)
		*(.data)
		_edata = .;
	}

	/* collect all uninitialized .bss sections */
	.bss (NOLOAD) : { 
		. = ALIGN(4);
		_sbss = .;
		*(.bss)
		_ebss = .;
	}
}
_romsize = _edata - _stext;
_sramsize = _ebss - _stext;
end = .;  /* define a global symbol marking the end of application */

