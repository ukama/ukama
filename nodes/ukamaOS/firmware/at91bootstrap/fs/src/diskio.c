/*-----------------------------------------------------------------------*/
/* Low level disk I/O module skeleton for FatFs     (C)ChaN, 2007        */
/*-----------------------------------------------------------------------*/
/* This is a stub disk I/O module that acts as front end of the existing */
/* disk I/O modules and attach it to FatFs module with common interface. */
/*-----------------------------------------------------------------------*/
#include "diskio.h"
#include "ffconf.h"
#include "integer.h"
#include "media.h"

//------------------------------------------------------------------------------
//         Internal variables

static volatile DSTATUS Stat = STA_NOINIT;	/* Disk status */

//------------------------------------------------------------------------------
/* Initialize a Drive                                                    */
/*-----------------------------------------------------------------------*/

DSTATUS disk_initialize(BYTE drv        /* Physical drive number (0..) */
    )
{
	if (drv) return STA_NOINIT;	
	
	if ((Stat & STA_NOINIT) != 0 && sdcard_initialize() == 0)
		Stat &= ~STA_NOINIT;

	return Stat;
}

/*-----------------------------------------------------------------------*/
/* Return Disk Status                                                    */
/*-----------------------------------------------------------------------*/

DSTATUS disk_status(BYTE drv    /* Physical drive number (0..) */
    )
{
	if (drv) return STA_NOINIT;

	return Stat;
}

/*-----------------------------------------------------------------------*/
/* Read Sector(s)                                                        */
/*-----------------------------------------------------------------------*/

DRESULT disk_read(BYTE drv,     /* Physical drive number (0..) */
                  BYTE *buff,  /* Data buffer to store read data */
                  DWORD sector, /* Start sector number (LBA) */
                  BYTE count    /* Sector count (1..255) */
    )
{
	if (drv || !count) return RES_PARERR;
	if (Stat & STA_NOINIT) return RES_NOTRDY;

	if (sdcard_block_read((unsigned int)sector,
				(unsigned int)count,
				(void *)buff) == count)
		return RES_OK;
	else
		return RES_ERROR;
}

/*-----------------------------------------------------------------------*/
/* Write Sector(s)                                                       */
/*-----------------------------------------------------------------------*/

#if _READONLY == 0
DRESULT disk_write(BYTE drv,    /* Physical drive number (0..) */
                   const BYTE * buff,   /* Data to be written */
                   DWORD sector,        /* Sector number (LBA) */
                   BYTE count   /* Sector count (1..255) */
    )
{
}
#endif  /* _READONLY */

/*-----------------------------------------------------------------------*/
/* Miscellaneous drive controls other than data read/write               */
/*-----------------------------------------------------------------------*/
#if _USE_IOCTL == 1
DRESULT disk_ioctl(BYTE drv,    /* Physical drive number (0..) */
                   BYTE ctrl,   /* Control code */
                   void *buff   /* Buffer to send/receive control data */
    )
{
}
#endif

/*----------------------------------------------------------------------*/
/* User Provided RTC Function for Fats module				*/
/*----------------------------------------------------------------------*/
/* Currnet time is returned with packed into a DWORD value.		*/
/* The bit field is as follows:						*/
/*   bit31:25  Year from 1980 (0..127)					*/
/*   bit24:21  Month (1..12)						*/
/*   bit20:16  Day in month(1..31)					*/
/*   bit15:11  Hour (0..23)						*/
/*   bit10:5   Minute (0..59)						*/
/*   bit4:0    Second / 2 (0..29)					*/

DWORD get_fattime(void)
{
    DWORD time;

    time = ((2009 - 1980) << 25)
        | (9 << 21)
        | (15 << 16)
        | (17 << 11)
        | (45 << 5)
        | ((59 / 2) << 0);

    return time;
}
