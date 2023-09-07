/* Copyright 2012 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Printf-like functionality for Chrome EC */

#include "printf.h"
#include "timer.h"
#include "util.h"

static const char error_str[] = "ERROR";

#define MAX_FORMAT 1024  /* Maximum chars in a single format field */

#ifndef CONFIG_DEBUG_PRINTF
static inline int divmod(uint64_t *n, int d)
{
	return uint64divmod(n, d);
}
#else /* CONFIG_DEBUG_PRINTF */
/* if we are optimizing for size, remove the 64-bit support */
#define NO_UINT64_SUPPORT
static inline int divmod(uint32_t *n, int d)
{
	int r = *n % d;
	*n /= d;
	return r;
}
#endif

/**
 * Convert the lowest nibble of a number to hex
 *
 * @param c	Number to extract lowest nibble from
 *
 * @return The corresponding ASCII character ('0' - 'f').
 */
static int hexdigit(int c)
{
	/* Strip off just the last nibble */
	c &= 0x0f;

	return c > 9 ? (c + 'a' - 10) : (c + '0');
}

/* Flags for vfnprintf() flags */
#define PF_LEFT		BIT(0)  /* Left-justify */
#define PF_PADZERO	BIT(1)  /* Pad with 0's not spaces */
#define PF_SIGN		BIT(2)  /* Add sign (+) for a positive number */
#define PF_64BIT	BIT(3)  /* Number is 64-bit */

int vfnprintf(int (*addchar)(void *context, int c), void *context,
	      const char *format, va_list args)
{
	/*
	 * Longest uint64 in decimal = 20
	 * Longest uint32 in binary  = 32
	 * + sign bit
	 * + terminating null
	 */
	char intbuf[34];
	int flags;
	int pad_width;
	int precision;
	char *vstr;
	int vlen;

	while (*format) {
		int c = *format++;
		char sign = 0;

		/* Copy normal characters */
		if (c != '%') {
			if (addchar(context, c))
				return EC_ERROR_OVERFLOW;
			continue;
		}

		/* Zero flags, now that we're in a format */
		flags = 0;

		/* Get first format character */
		c = *format++;

		/* Send "%" for "%%" input */
		if (c == '%' || c == '\0') {
			if (addchar(context, '%'))
				return EC_ERROR_OVERFLOW;

			if (c == '\0')
				break;

			continue;
		}

		/* Handle %c */
		if (c == 'c') {
			c = va_arg(args, int);
			if (addchar(context, c))
				return EC_ERROR_OVERFLOW;
			continue;
		}

		/* Handle left-justification ("%-5s") */
		if (c == '-') {
			flags |= PF_LEFT;
			c = *format++;
		}

		/* Handle positive sign (%+d) */
		if (c == '+') {
			flags |= PF_SIGN;
			c = *format++;
		}

		/* Handle padding with 0's */
		if (c == '0') {
			flags |= PF_PADZERO;
			c = *format++;
		}

		/* Count padding length */
		pad_width = 0;
		if (c == '*') {
			pad_width = va_arg(args, int);
			c = *format++;
		} else {
			while (c >= '0' && c <= '9') {
				pad_width = (10 * pad_width) + c - '0';
				c = *format++;
			}
		}
		if (pad_width < 0 || pad_width > MAX_FORMAT) {
			/* Sanity check for precision failed */
			format = error_str;
			continue;
		}

		/* Count precision */
		precision = -1;
		if (c == '.') {
			c = *format++;
			if (c == '*') {
				precision = va_arg(args, int);
				c = *format++;
			} else {
				precision = 0;
				while (c >= '0' && c <= '9') {
					precision = (10 * precision) + c - '0';
					c = *format++;
				}
			}
			if (precision < 0 || precision > MAX_FORMAT) {
				/* Sanity check for precision failed */
				format = error_str;
				continue;
			}
		}

		if (c == 's') {
			vstr = va_arg(args, char *);
			if (vstr == NULL)
				vstr = "(NULL)";
		} else if (c == 'h') {
			/* Hex dump output */
			vstr = va_arg(args, char *);

			if (precision < 0) {
				/* Hex dump requires precision */
				format = error_str;
				continue;
			}

			/*
			 * Divide pad_width instead of multiplying
			 * precision to avoid overflow error
			 * in the condition.
			 * The "/2" and "2*" can be optimized by
			 * the compiler.
			 */
			if ((pad_width/2) >= precision)
				pad_width -= 2*precision;
			else
				pad_width = 0;

			while (pad_width > 0 && !(flags & PF_LEFT)) {
				if (addchar(context,
					    flags & PF_PADZERO ? '0' : ' '))
					return EC_ERROR_OVERFLOW;
				pad_width--;
			}
			for (; precision; precision--, vstr++) {
				if (addchar(context, hexdigit(*vstr >> 4)) ||
				    addchar(context, hexdigit(*vstr)))
					return EC_ERROR_OVERFLOW;
			}
			while (pad_width > 0 && (flags & PF_LEFT)) {
				if (addchar(context, ' '))
					return EC_ERROR_OVERFLOW;
				pad_width--;
			}

			continue;
		} else {
			int base = 10;
#ifdef NO_UINT64_SUPPORT
			uint32_t v;

			v = va_arg(args, uint32_t);
#else /* NO_UINT64_SUPPORT */
			uint64_t v;

			/* Handle length */
			if (c == 'l') {
				flags |= PF_64BIT;
				c = *format++;
			}

			/* Special-case: %T = current time */
			if (c == 'T') {
				v = get_time().val;
				flags |= PF_64BIT;
#ifdef CONFIG_CONSOLE_VERBOSE
				precision = 6;
#else
				precision = 3;
				v /= 1000;
#endif
			} else if (flags & PF_64BIT) {
				v = va_arg(args, uint64_t);
			} else {
				v = va_arg(args, uint32_t);
			}
#endif

			switch (c) {
#ifdef CONFIG_PRINTF_LEGACY_LI_FORMAT
			case 'i':
				/* force 32-bit for compatibility */
				flags &= ~PF_64BIT;
				/* fall-through */
#endif /* CONFIG_PRINTF_LEGACY_LI_FORMAT */
			case 'd':
				if (flags & PF_64BIT) {
					if ((int64_t)v < 0) {
						sign = '-';
						if (v != (1ULL << 63))
							v = -v;
					} else if (flags & PF_SIGN) {
						sign = '+';
					}
				} else {
					if ((int)v < 0) {
						sign = '-';
						if (v != (1ULL << 31))
							v = -(int)v;
					} else if (flags & PF_SIGN) {
						sign = '+';
					}
				}
				break;
			case 'u':
			case 'T':
				break;
			case 'X':
			case 'x':
			case 'p':
				base = 16;
				break;
			case 'b':
				base = 2;
				break;
			default:
				format = error_str;
			}
			if (format == error_str)
				continue; /* Bad format specifier */

			/*
			 * Convert integer to string, starting at end of
			 * buffer and working backwards.
			 */
			vstr = intbuf + sizeof(intbuf) - 1;
			*(vstr) = '\0';

			/*
			 * Fixed-point precision must fit in our buffer.
			 * Leave space for "0." and the terminating null.
			 */
			if (precision > (int)(sizeof(intbuf) - 3))
				precision = sizeof(intbuf) - 3;

			/*
			 * Handle digits to right of decimal for fixed point
			 * numbers.
			 */
			for (vlen = 0; vlen < precision; vlen++)
				*(--vstr) = '0' + divmod(&v, 10);
			if (precision >= 0)
				*(--vstr) = '.';

			if (!v)
				*(--vstr) = '0';

			while (v) {
				int digit = divmod(&v, base);
				if (digit < 10)
					*(--vstr) = '0' + digit;
				else if (c == 'X')
					*(--vstr) = 'A' + digit - 10;
				else
					*(--vstr) = 'a' + digit - 10;
			}

			if (sign)
				*(--vstr) = sign;

			/*
			 * Precision field was interpreted by fixed-point
			 * logic, so clear it.
			 */
			precision = -1;
		}

		/* No padding strings to wider than the precision */
		if (precision >= 0 && pad_width > precision)
			pad_width = precision;

		if (precision < 0) {
			/* If precision is unset, print everything */
			vlen = strlen(vstr);
			precision = MAX(vlen, pad_width);
		} else {
			/*
			 * If precision is set, ensure that we do not
			 * overrun it
			 */
			vlen = strnlen(vstr, precision);
		}


		while (vlen < pad_width && !(flags & PF_LEFT)) {
			if (addchar(context, flags & PF_PADZERO ? '0' : ' '))
				return EC_ERROR_OVERFLOW;
			vlen++;
		}
		while (--precision >= 0 && *vstr)
			if (addchar(context, *vstr++))
				return EC_ERROR_OVERFLOW;
		while (vlen < pad_width && flags & PF_LEFT) {
			if (addchar(context, ' '))
				return EC_ERROR_OVERFLOW;
			vlen++;
		}
	}

	/* If we're still here, we consumed all output */
	return EC_SUCCESS;
}

/* Context for snprintf() */
struct snprintf_context {
	char *str;
	int size;
};

/**
 * Add a character to the string context.
 *
 * @param context	Context receiving character
 * @param c		Character to add
 * @return 0 if character added, 1 if character dropped because no space.
 */
static int snprintf_addchar(void *context, int c)
{
	struct snprintf_context *ctx = (struct snprintf_context *)context;

	if (!ctx->size)
		return 1;

	*(ctx->str++) = c;
	ctx->size--;
	return 0;
}

int snprintf(char *str, int size, const char *format, ...)
{
	va_list args;
	int rv;

	va_start(args, format);
	rv = vsnprintf(str, size, format, args);
	va_end(args);

	return rv;
}

int vsnprintf(char *str, int size, const char *format, va_list args)
{
	struct snprintf_context ctx;
	int rv;

	if (!str || !format || size <= 0)
		return -EC_ERROR_INVAL;

	ctx.str = str;
	ctx.size = size - 1;  /* Reserve space for terminating '\0' */

	rv = vfnprintf(snprintf_addchar, &ctx, format, args);

	/* Terminate string */
	*ctx.str = '\0';

	return (rv == EC_SUCCESS) ? (ctx.str - str) : -rv;
}
