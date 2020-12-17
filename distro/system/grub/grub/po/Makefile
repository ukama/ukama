# Makefile for PO directory in any package using GNU gettext.
# Copyright (C) 1995-1997, 2000-2007, 2009-2010 by Ulrich Drepper <drepper@gnu.ai.mit.edu>
#
# This file can be copied and used freely without restrictions.  It can
# be used in projects which are not available under the GNU General Public
# License but which still want to provide support for the GNU gettext
# functionality.
# Please note that the actual code of GNU gettext is covered by the GNU
# General Public License and is *not* in the public domain.
#
# Origin: gettext-0.18.3
GETTEXT_MACRO_VERSION = 0.18

PACKAGE = grub
VERSION = 2.04
PACKAGE_BUGREPORT = bug-grub@gnu.org

SED = /bin/sed
SHELL = /bin/sh


srcdir = .
top_srcdir = ..


prefix = /ukamaOS/_ukamafs
exec_prefix = /ukamaOS/_ukamafs
datarootdir = /ukamaOS/_ukamafs/share
datadir = ${datarootdir}
localedir = ${datarootdir}/locale
gettextsrcdir = $(datadir)/gettext/po

INSTALL = /usr/bin/install -c
INSTALL_DATA = ${INSTALL} -m 644

# We use $(mkdir_p).
# In automake <= 1.9.x, $(mkdir_p) is defined either as "mkdir -p --" or as
# "$(mkinstalldirs)" or as "$(install_sh) -d". For these automake versions,
# ${SHELL} /ukamaOS/distro/system/grub/grub/build-aux/install-sh does not start with $(SHELL), so we add it.
# In automake >= 1.10, $(MKDIR_P) is derived from ${MKDIR_P}, which is defined
# either as "/path/to/mkdir -p" or ".../install-sh -c -d". For these automake
# versions, $(mkinstalldirs) and $(install_sh) are unused.
mkinstalldirs = $(SHELL) ${SHELL} /ukamaOS/distro/system/grub/grub/build-aux/install-sh -d
install_sh = $(SHELL) ${SHELL} /ukamaOS/distro/system/grub/grub/build-aux/install-sh
MKDIR_P = /bin/mkdir -p
mkdir_p = $(MKDIR_P)

GMSGFMT_ = /usr/bin/msgfmt
GMSGFMT_no = /usr/bin/msgfmt
GMSGFMT_yes = /usr/bin/msgfmt
GMSGFMT = $(GMSGFMT_$(USE_MSGCTXT))
MSGFMT_ = /usr/bin/msgfmt
MSGFMT_no = /usr/bin/msgfmt
MSGFMT_yes = /usr/bin/msgfmt
MSGFMT = $(MSGFMT_$(USE_MSGCTXT))
XGETTEXT_ = /usr/bin/xgettext
XGETTEXT_no = /usr/bin/xgettext
XGETTEXT_yes = /usr/bin/xgettext
XGETTEXT = $(XGETTEXT_$(USE_MSGCTXT))
MSGMERGE = msgmerge
MSGMERGE_UPDATE = /usr/bin/msgmerge --update
MSGINIT = msginit
MSGCONV = msgconv
MSGFILTER = msgfilter

POFILES =  ast.po ca.po da.po de.po de@hebrew.po de_CH.po en@arabic.po en@cyrillic.po en@greek.po en@hebrew.po en@piglatin.po en@quot.po eo.po es.po fi.po fr.po gl.po hr.po hu.po id.po it.po ja.po ko.po lt.po nb.po nl.po pa.po pl.po pt.po pt_BR.po ro.po ru.po sl.po sr.po sv.po tr.po uk.po vi.po zh_CN.po zh_TW.po
GMOFILES =  ast.gmo ca.gmo da.gmo de.gmo de@hebrew.gmo de_CH.gmo en@arabic.gmo en@cyrillic.gmo en@greek.gmo en@hebrew.gmo en@piglatin.gmo en@quot.gmo eo.gmo es.gmo fi.gmo fr.gmo gl.gmo hr.gmo hu.gmo id.gmo it.gmo ja.gmo ko.gmo lt.gmo nb.gmo nl.gmo pa.gmo pl.gmo pt.gmo pt_BR.gmo ro.gmo ru.gmo sl.gmo sr.gmo sv.gmo tr.gmo uk.gmo vi.gmo zh_CN.gmo zh_TW.gmo
UPDATEPOFILES =  ast.po-update ca.po-update da.po-update de.po-update de@hebrew.po-update de_CH.po-update en@arabic.po-update en@cyrillic.po-update en@greek.po-update en@hebrew.po-update en@piglatin.po-update en@quot.po-update eo.po-update es.po-update fi.po-update fr.po-update gl.po-update hr.po-update hu.po-update id.po-update it.po-update ja.po-update ko.po-update lt.po-update nb.po-update nl.po-update pa.po-update pl.po-update pt.po-update pt_BR.po-update ro.po-update ru.po-update sl.po-update sr.po-update sv.po-update tr.po-update uk.po-update vi.po-update zh_CN.po-update zh_TW.po-update
DUMMYPOFILES =  ast.nop ca.nop da.nop de.nop de@hebrew.nop de_CH.nop en@arabic.nop en@cyrillic.nop en@greek.nop en@hebrew.nop en@piglatin.nop en@quot.nop eo.nop es.nop fi.nop fr.nop gl.nop hr.nop hu.nop id.nop it.nop ja.nop ko.nop lt.nop nb.nop nl.nop pa.nop pl.nop pt.nop pt_BR.nop ro.nop ru.nop sl.nop sr.nop sv.nop tr.nop uk.nop vi.nop zh_CN.nop zh_TW.nop
DISTFILES.common = Makefile.in.in remove-potcdate.sin \
$(DISTFILES.common.extra1) $(DISTFILES.common.extra2) $(DISTFILES.common.extra3)
DISTFILES = $(DISTFILES.common) Makevars POTFILES.in \
$(POFILES) $(GMOFILES) \
$(DISTFILES.extra1) $(DISTFILES.extra2) $(DISTFILES.extra3)

POTFILES = \
     .././grub-core/boot/decompressor/minilib.c \
     .././grub-core/boot/decompressor/none.c \
     .././grub-core/boot/decompressor/xz.c \
     .././grub-core/bus/bonito.c \
     .././grub-core/bus/cs5536.c \
     .././grub-core/bus/emu/pci.c \
     .././grub-core/bus/fdt.c \
     .././grub-core/bus/i386/ieee1275/pci.c \
     .././grub-core/bus/pci.c \
     .././grub-core/bus/spi/rk3288_spi.c \
     .././grub-core/bus/usb/ehci-fdt.c \
     .././grub-core/bus/usb/ehci-pci.c \
     .././grub-core/bus/usb/ehci.c \
     .././grub-core/bus/usb/ohci.c \
     .././grub-core/bus/usb/serial/common.c \
     .././grub-core/bus/usb/serial/ftdi.c \
     .././grub-core/bus/usb/serial/pl2303.c \
     .././grub-core/bus/usb/serial/usbdebug_late.c \
     .././grub-core/bus/usb/uhci.c \
     .././grub-core/bus/usb/usb.c \
     .././grub-core/bus/usb/usbhub.c \
     .././grub-core/bus/usb/usbtrans.c \
     .././grub-core/commands/acpi.c \
     .././grub-core/commands/acpihalt.c \
     .././grub-core/commands/arc/lsdev.c \
     .././grub-core/commands/blocklist.c \
     .././grub-core/commands/boot.c \
     .././grub-core/commands/boottime.c \
     .././grub-core/commands/cacheinfo.c \
     .././grub-core/commands/cat.c \
     .././grub-core/commands/cmp.c \
     .././grub-core/commands/configfile.c \
     .././grub-core/commands/date.c \
     .././grub-core/commands/echo.c \
     .././grub-core/commands/efi/efifwsetup.c \
     .././grub-core/commands/efi/fixvideo.c \
     .././grub-core/commands/efi/loadbios.c \
     .././grub-core/commands/efi/lsefi.c \
     .././grub-core/commands/efi/lsefimmap.c \
     .././grub-core/commands/efi/lsefisystab.c \
     .././grub-core/commands/efi/lssal.c \
     .././grub-core/commands/efi/shim_lock.c \
     .././grub-core/commands/efi/tpm.c \
     .././grub-core/commands/eval.c \
     .././grub-core/commands/extcmd.c \
     .././grub-core/commands/file.c \
     .././grub-core/commands/file32.c \
     .././grub-core/commands/file64.c \
     .././grub-core/commands/fileXX.c \
     .././grub-core/commands/gptsync.c \
     .././grub-core/commands/halt.c \
     .././grub-core/commands/hashsum.c \
     .././grub-core/commands/hdparm.c \
     .././grub-core/commands/help.c \
     .././grub-core/commands/hexdump.c \
     .././grub-core/commands/i386/cmosdump.c \
     .././grub-core/commands/i386/cmostest.c \
     .././grub-core/commands/i386/coreboot/cb_timestamps.c \
     .././grub-core/commands/i386/coreboot/cbls.c \
     .././grub-core/commands/i386/cpuid.c \
     .././grub-core/commands/i386/pc/drivemap.c \
     .././grub-core/commands/i386/pc/halt.c \
     .././grub-core/commands/i386/pc/lsapm.c \
     .././grub-core/commands/i386/pc/play.c \
     .././grub-core/commands/i386/pc/sendkey.c \
     .././grub-core/commands/i386/rdmsr.c \
     .././grub-core/commands/i386/wrmsr.c \
     .././grub-core/commands/ieee1275/suspend.c \
     .././grub-core/commands/iorw.c \
     .././grub-core/commands/keylayouts.c \
     .././grub-core/commands/keystatus.c \
     .././grub-core/commands/legacycfg.c \
     .././grub-core/commands/loadenv.c \
     .././grub-core/commands/ls.c \
     .././grub-core/commands/lsacpi.c \
     .././grub-core/commands/lsmmap.c \
     .././grub-core/commands/lspci.c \
     .././grub-core/commands/macbless.c \
     .././grub-core/commands/memrw.c \
     .././grub-core/commands/menuentry.c \
     .././grub-core/commands/minicmd.c \
     .././grub-core/commands/mips/loongson/lsspd.c \
     .././grub-core/commands/nativedisk.c \
     .././grub-core/commands/parttool.c \
     .././grub-core/commands/password.c \
     .././grub-core/commands/password_pbkdf2.c \
     .././grub-core/commands/pcidump.c \
     .././grub-core/commands/pgp.c \
     .././grub-core/commands/probe.c \
     .././grub-core/commands/read.c \
     .././grub-core/commands/reboot.c \
     .././grub-core/commands/regexp.c \
     .././grub-core/commands/search.c \
     .././grub-core/commands/search_file.c \
     .././grub-core/commands/search_label.c \
     .././grub-core/commands/search_uuid.c \
     .././grub-core/commands/search_wrap.c \
     .././grub-core/commands/setpci.c \
     .././grub-core/commands/sleep.c \
     .././grub-core/commands/syslinuxcfg.c \
     .././grub-core/commands/terminal.c \
     .././grub-core/commands/test.c \
     .././grub-core/commands/testload.c \
     .././grub-core/commands/testspeed.c \
     .././grub-core/commands/time.c \
     .././grub-core/commands/tpm.c \
     .././grub-core/commands/tr.c \
     .././grub-core/commands/true.c \
     .././grub-core/commands/usbtest.c \
     .././grub-core/commands/verifiers.c \
     .././grub-core/commands/videoinfo.c \
     .././grub-core/commands/videotest.c \
     .././grub-core/commands/wildcard.c \
     .././grub-core/commands/xen/lsxen.c \
     .././grub-core/commands/xnu_uuid.c \
     .././grub-core/disk/AFSplitter.c \
     .././grub-core/disk/ahci.c \
     .././grub-core/disk/arc/arcdisk.c \
     .././grub-core/disk/ata.c \
     .././grub-core/disk/cryptodisk.c \
     .././grub-core/disk/diskfilter.c \
     .././grub-core/disk/dmraid_nvidia.c \
     .././grub-core/disk/efi/efidisk.c \
     .././grub-core/disk/geli.c \
     .././grub-core/disk/host.c \
     .././grub-core/disk/i386/pc/biosdisk.c \
     .././grub-core/disk/ieee1275/nand.c \
     .././grub-core/disk/ieee1275/obdisk.c \
     .././grub-core/disk/ieee1275/ofdisk.c \
     .././grub-core/disk/ldm.c \
     .././grub-core/disk/loopback.c \
     .././grub-core/disk/luks.c \
     .././grub-core/disk/lvm.c \
     .././grub-core/disk/mdraid1x_linux.c \
     .././grub-core/disk/mdraid_linux.c \
     .././grub-core/disk/mdraid_linux_be.c \
     .././grub-core/disk/memdisk.c \
     .././grub-core/disk/pata.c \
     .././grub-core/disk/raid5_recover.c \
     .././grub-core/disk/raid6_recover.c \
     .././grub-core/disk/scsi.c \
     .././grub-core/disk/uboot/ubootdisk.c \
     .././grub-core/disk/usbms.c \
     .././grub-core/disk/xen/xendisk.c \
     .././grub-core/efiemu/i386/coredetect.c \
     .././grub-core/efiemu/i386/loadcore32.c \
     .././grub-core/efiemu/i386/loadcore64.c \
     .././grub-core/efiemu/i386/nocfgtables.c \
     .././grub-core/efiemu/i386/pc/cfgtables.c \
     .././grub-core/efiemu/loadcore.c \
     .././grub-core/efiemu/loadcore32.c \
     .././grub-core/efiemu/loadcore64.c \
     .././grub-core/efiemu/loadcore_common.c \
     .././grub-core/efiemu/main.c \
     .././grub-core/efiemu/mm.c \
     .././grub-core/efiemu/pnvram.c \
     .././grub-core/efiemu/prepare.c \
     .././grub-core/efiemu/prepare32.c \
     .././grub-core/efiemu/prepare64.c \
     .././grub-core/efiemu/runtime/config.h \
     .././grub-core/efiemu/runtime/efiemu.c \
     .././grub-core/efiemu/symbols.c \
     .././grub-core/font/font.c \
     .././grub-core/font/font_cmd.c \
     .././grub-core/fs/affs.c \
     .././grub-core/fs/afs.c \
     .././grub-core/fs/archelp.c \
     .././grub-core/fs/bfs.c \
     .././grub-core/fs/btrfs.c \
     .././grub-core/fs/cbfs.c \
     .././grub-core/fs/cpio.c \
     .././grub-core/fs/cpio_be.c \
     .././grub-core/fs/cpio_common.c \
     .././grub-core/fs/exfat.c \
     .././grub-core/fs/ext2.c \
     .././grub-core/fs/f2fs.c \
     .././grub-core/fs/fat.c \
     .././grub-core/fs/fshelp.c \
     .././grub-core/fs/hfs.c \
     .././grub-core/fs/hfsplus.c \
     .././grub-core/fs/hfspluscomp.c \
     .././grub-core/fs/iso9660.c \
     .././grub-core/fs/jfs.c \
     .././grub-core/fs/minix.c \
     .././grub-core/fs/minix2.c \
     .././grub-core/fs/minix2_be.c \
     .././grub-core/fs/minix3.c \
     .././grub-core/fs/minix3_be.c \
     .././grub-core/fs/minix_be.c \
     .././grub-core/fs/newc.c \
     .././grub-core/fs/nilfs2.c \
     .././grub-core/fs/ntfs.c \
     .././grub-core/fs/ntfscomp.c \
     .././grub-core/fs/odc.c \
     .././grub-core/fs/proc.c \
     .././grub-core/fs/reiserfs.c \
     .././grub-core/fs/romfs.c \
     .././grub-core/fs/sfs.c \
     .././grub-core/fs/squash4.c \
     .././grub-core/fs/tar.c \
     .././grub-core/fs/udf.c \
     .././grub-core/fs/ufs.c \
     .././grub-core/fs/ufs2.c \
     .././grub-core/fs/ufs_be.c \
     .././grub-core/fs/xfs.c \
     .././grub-core/fs/zfs/zfs.c \
     .././grub-core/fs/zfs/zfs_fletcher.c \
     .././grub-core/fs/zfs/zfs_lz4.c \
     .././grub-core/fs/zfs/zfs_lzjb.c \
     .././grub-core/fs/zfs/zfs_sha256.c \
     .././grub-core/fs/zfs/zfscrypt.c \
     .././grub-core/fs/zfs/zfsinfo.c \
     .././grub-core/gdb/cstub.c \
     .././grub-core/gdb/gdb.c \
     .././grub-core/gdb/i386/idt.c \
     .././grub-core/gdb/i386/signal.c \
     .././grub-core/gentrigtables.c \
     .././grub-core/gettext/gettext.c \
     .././grub-core/gfxmenu/font.c \
     .././grub-core/gfxmenu/gfxmenu.c \
     .././grub-core/gfxmenu/gui_box.c \
     .././grub-core/gfxmenu/gui_canvas.c \
     .././grub-core/gfxmenu/gui_circular_progress.c \
     .././grub-core/gfxmenu/gui_image.c \
     .././grub-core/gfxmenu/gui_label.c \
     .././grub-core/gfxmenu/gui_list.c \
     .././grub-core/gfxmenu/gui_progress_bar.c \
     .././grub-core/gfxmenu/gui_string_util.c \
     .././grub-core/gfxmenu/gui_util.c \
     .././grub-core/gfxmenu/icon_manager.c \
     .././grub-core/gfxmenu/theme_loader.c \
     .././grub-core/gfxmenu/view.c \
     .././grub-core/gfxmenu/widget-box.c \
     .././grub-core/hello/hello.c \
     .././grub-core/hook/datehook.c \
     .././grub-core/io/bufio.c \
     .././grub-core/io/gzio.c \
     .././grub-core/io/lzopio.c \
     .././grub-core/io/offset.c \
     .././grub-core/io/xzio.c \
     .././grub-core/kern/acpi.c \
     .././grub-core/kern/arm/cache.c \
     .././grub-core/kern/arm/coreboot/cbtable.c \
     .././grub-core/kern/arm/coreboot/dma.c \
     .././grub-core/kern/arm/coreboot/init.c \
     .././grub-core/kern/arm/coreboot/timer.c \
     .././grub-core/kern/arm/dl.c \
     .././grub-core/kern/arm/dl_helper.c \
     .././grub-core/kern/arm/efi/init.c \
     .././grub-core/kern/arm/uboot/init.c \
     .././grub-core/kern/arm64/cache.c \
     .././grub-core/kern/arm64/dl.c \
     .././grub-core/kern/arm64/dl_helper.c \
     .././grub-core/kern/arm64/efi/init.c \
     .././grub-core/kern/command.c \
     .././grub-core/kern/compiler-rt.c \
     .././grub-core/kern/coreboot/cbtable.c \
     .././grub-core/kern/coreboot/mmap.c \
     .././grub-core/kern/corecmd.c \
     .././grub-core/kern/device.c \
     .././grub-core/kern/disk.c \
     .././grub-core/kern/disk_common.c \
     .././grub-core/kern/dl.c \
     .././grub-core/kern/efi/acpi.c \
     .././grub-core/kern/efi/efi.c \
     .././grub-core/kern/efi/fdt.c \
     .././grub-core/kern/efi/init.c \
     .././grub-core/kern/efi/mm.c \
     .././grub-core/kern/elf.c \
     .././grub-core/kern/elfXX.c \
     .././grub-core/kern/emu/argp_common.c \
     .././grub-core/kern/emu/cache.c \
     .././grub-core/kern/emu/full.c \
     .././grub-core/kern/emu/hostdisk.c \
     .././grub-core/kern/emu/hostfs.c \
     .././grub-core/kern/emu/lite.c \
     .././grub-core/kern/emu/main.c \
     .././grub-core/kern/emu/misc.c \
     .././grub-core/kern/emu/mm.c \
     .././grub-core/kern/emu/time.c \
     .././grub-core/kern/env.c \
     .././grub-core/kern/err.c \
     .././grub-core/kern/file.c \
     .././grub-core/kern/fs.c \
     .././grub-core/kern/generic/millisleep.c \
     .././grub-core/kern/generic/rtc_get_time_ms.c \
     .././grub-core/kern/i386/coreboot/cbtable.c \
     .././grub-core/kern/i386/coreboot/init.c \
     .././grub-core/kern/i386/dl.c \
     .././grub-core/kern/i386/efi/init.c \
     .././grub-core/kern/i386/efi/tsc.c \
     .././grub-core/kern/i386/multiboot_mmap.c \
     .././grub-core/kern/i386/pc/acpi.c \
     .././grub-core/kern/i386/pc/init.c \
     .././grub-core/kern/i386/pc/mmap.c \
     .././grub-core/kern/i386/qemu/init.c \
     .././grub-core/kern/i386/qemu/mmap.c \
     .././grub-core/kern/i386/tsc.c \
     .././grub-core/kern/i386/tsc_pit.c \
     .././grub-core/kern/i386/tsc_pmtimer.c \
     .././grub-core/kern/i386/xen/pvh.c \
     .././grub-core/kern/i386/xen/tsc.c \
     .././grub-core/kern/ia64/cache.c \
     .././grub-core/kern/ia64/dl.c \
     .././grub-core/kern/ia64/dl_helper.c \
     .././grub-core/kern/ia64/efi/init.c \
     .././grub-core/kern/ieee1275/cmain.c \
     .././grub-core/kern/ieee1275/ieee1275.c \
     .././grub-core/kern/ieee1275/init.c \
     .././grub-core/kern/ieee1275/mmap.c \
     .././grub-core/kern/ieee1275/openfw.c \
     .././grub-core/kern/list.c \
     .././grub-core/kern/main.c \
     .././grub-core/kern/mips/arc/init.c \
     .././grub-core/kern/mips/dl.c \
     .././grub-core/kern/mips/init.c \
     .././grub-core/kern/mips/loongson/init.c \
     .././grub-core/kern/mips/qemu_mips/init.c \
     .././grub-core/kern/misc.c \
     .././grub-core/kern/mm.c \
     .././grub-core/kern/parser.c \
     .././grub-core/kern/partition.c \
     .././grub-core/kern/powerpc/dl.c \
     .././grub-core/kern/rescue_parser.c \
     .././grub-core/kern/rescue_reader.c \
     .././grub-core/kern/riscv/cache.c \
     .././grub-core/kern/riscv/dl.c \
     .././grub-core/kern/riscv/efi/init.c \
     .././grub-core/kern/sparc64/dl.c \
     .././grub-core/kern/sparc64/ieee1275/ieee1275.c \
     .././grub-core/kern/term.c \
     .././grub-core/kern/time.c \
     .././grub-core/kern/uboot/hw.c \
     .././grub-core/kern/uboot/init.c \
     .././grub-core/kern/uboot/uboot.c \
     .././grub-core/kern/vga_init.c \
     .././grub-core/kern/x86_64/dl.c \
     .././grub-core/kern/xen/init.c \
     .././grub-core/lib/LzFind.c \
     .././grub-core/lib/LzmaDec.c \
     .././grub-core/lib/LzmaEnc.c \
     .././grub-core/lib/adler32.c \
     .././grub-core/lib/arc/datetime.c \
     .././grub-core/lib/arg.c \
     .././grub-core/lib/backtrace.c \
     .././grub-core/lib/cmdline.c \
     .././grub-core/lib/cmos_datetime.c \
     .././grub-core/lib/crc.c \
     .././grub-core/lib/crc64.c \
     .././grub-core/lib/crypto.c \
     .././grub-core/lib/disk.c \
     .././grub-core/lib/division.c \
     .././grub-core/lib/dummy/datetime.c \
     .././grub-core/lib/dummy/halt.c \
     .././grub-core/lib/dummy/reboot.c \
     .././grub-core/lib/efi/datetime.c \
     .././grub-core/lib/efi/halt.c \
     .././grub-core/lib/efi/relocator.c \
     .././grub-core/lib/emu/halt.c \
     .././grub-core/lib/envblk.c \
     .././grub-core/lib/fake_module.c \
     .././grub-core/lib/fdt.c \
     .././grub-core/lib/getline.c \
     .././grub-core/lib/gnulib/_Noreturn.h \
     .././grub-core/lib/gnulib/alloca.c \
     .././grub-core/lib/gnulib/alloca.in.h \
     .././grub-core/lib/gnulib/arg-nonnull.h \
     .././grub-core/lib/gnulib/argp-ba.c \
     .././grub-core/lib/gnulib/argp-eexst.c \
     .././grub-core/lib/gnulib/argp-fmtstream.c \
     .././grub-core/lib/gnulib/argp-fmtstream.h \
     .././grub-core/lib/gnulib/argp-fs-xinl.c \
     .././grub-core/lib/gnulib/argp-help.c \
     .././grub-core/lib/gnulib/argp-namefrob.h \
     .././grub-core/lib/gnulib/argp-parse.c \
     .././grub-core/lib/gnulib/argp-pin.c \
     .././grub-core/lib/gnulib/argp-pv.c \
     .././grub-core/lib/gnulib/argp-pvh.c \
     .././grub-core/lib/gnulib/argp-xinl.c \
     .././grub-core/lib/gnulib/argp.h \
     .././grub-core/lib/gnulib/asnprintf.c \
     .././grub-core/lib/gnulib/assure.h \
     .././grub-core/lib/gnulib/basename-lgpl.c \
     .././grub-core/lib/gnulib/btowc.c \
     .././grub-core/lib/gnulib/c++defs.h \
     .././grub-core/lib/gnulib/cdefs.h \
     .././grub-core/lib/gnulib/chdir-long.c \
     .././grub-core/lib/gnulib/chdir-long.h \
     .././grub-core/lib/gnulib/cloexec.c \
     .././grub-core/lib/gnulib/cloexec.h \
     .././grub-core/lib/gnulib/close.c \
     .././grub-core/lib/gnulib/dirent.in.h \
     .././grub-core/lib/gnulib/dirfd.c \
     .././grub-core/lib/gnulib/dirname-lgpl.c \
     .././grub-core/lib/gnulib/dirname.h \
     .././grub-core/lib/gnulib/dosname.h \
     .././grub-core/lib/gnulib/dup-safer-flag.c \
     .././grub-core/lib/gnulib/dup-safer.c \
     .././grub-core/lib/gnulib/dup2.c \
     .././grub-core/lib/gnulib/errno.in.h \
     .././grub-core/lib/gnulib/error.c \
     .././grub-core/lib/gnulib/error.h \
     .././grub-core/lib/gnulib/exitfail.c \
     .././grub-core/lib/gnulib/exitfail.h \
     .././grub-core/lib/gnulib/fchdir.c \
     .././grub-core/lib/gnulib/fcntl.c \
     .././grub-core/lib/gnulib/fcntl.in.h \
     .././grub-core/lib/gnulib/fd-hook.c \
     .././grub-core/lib/gnulib/fd-hook.h \
     .././grub-core/lib/gnulib/fd-safer-flag.c \
     .././grub-core/lib/gnulib/fd-safer.c \
     .././grub-core/lib/gnulib/filename.h \
     .././grub-core/lib/gnulib/filenamecat-lgpl.c \
     .././grub-core/lib/gnulib/filenamecat.h \
     .././grub-core/lib/gnulib/flexmember.h \
     .././grub-core/lib/gnulib/float+.h \
     .././grub-core/lib/gnulib/float.c \
     .././grub-core/lib/gnulib/float.in.h \
     .././grub-core/lib/gnulib/fnmatch.c \
     .././grub-core/lib/gnulib/fnmatch.in.h \
     .././grub-core/lib/gnulib/fnmatch_loop.c \
     .././grub-core/lib/gnulib/fstat.c \
     .././grub-core/lib/gnulib/getcwd-lgpl.c \
     .././grub-core/lib/gnulib/getdelim.c \
     .././grub-core/lib/gnulib/getdtablesize.c \
     .././grub-core/lib/gnulib/getline.c \
     .././grub-core/lib/gnulib/getopt-cdefs.in.h \
     .././grub-core/lib/gnulib/getopt-core.h \
     .././grub-core/lib/gnulib/getopt-ext.h \
     .././grub-core/lib/gnulib/getopt-pfx-core.h \
     .././grub-core/lib/gnulib/getopt-pfx-ext.h \
     .././grub-core/lib/gnulib/getopt.c \
     .././grub-core/lib/gnulib/getopt.in.h \
     .././grub-core/lib/gnulib/getopt1.c \
     .././grub-core/lib/gnulib/getopt_int.h \
     .././grub-core/lib/gnulib/getprogname.c \
     .././grub-core/lib/gnulib/getprogname.h \
     .././grub-core/lib/gnulib/gettext.h \
     .././grub-core/lib/gnulib/glthread/lock.c \
     .././grub-core/lib/gnulib/glthread/lock.h \
     .././grub-core/lib/gnulib/glthread/threadlib.c \
     .././grub-core/lib/gnulib/hard-locale.c \
     .././grub-core/lib/gnulib/hard-locale.h \
     .././grub-core/lib/gnulib/intprops.h \
     .././grub-core/lib/gnulib/itold.c \
     .././grub-core/lib/gnulib/langinfo.in.h \
     .././grub-core/lib/gnulib/libc-config.h \
     .././grub-core/lib/gnulib/limits.in.h \
     .././grub-core/lib/gnulib/localcharset.c \
     .././grub-core/lib/gnulib/localcharset.h \
     .././grub-core/lib/gnulib/locale.in.h \
     .././grub-core/lib/gnulib/localeconv.c \
     .././grub-core/lib/gnulib/malloc.c \
     .././grub-core/lib/gnulib/malloca.c \
     .././grub-core/lib/gnulib/malloca.h \
     .././grub-core/lib/gnulib/mbrtowc.c \
     .././grub-core/lib/gnulib/mbsinit.c \
     .././grub-core/lib/gnulib/mbsrtowcs-impl.h \
     .././grub-core/lib/gnulib/mbsrtowcs-state.c \
     .././grub-core/lib/gnulib/mbsrtowcs.c \
     .././grub-core/lib/gnulib/mbswidth.c \
     .././grub-core/lib/gnulib/mbswidth.h \
     .././grub-core/lib/gnulib/mbtowc-impl.h \
     .././grub-core/lib/gnulib/mbtowc.c \
     .././grub-core/lib/gnulib/memchr.c \
     .././grub-core/lib/gnulib/mempcpy.c \
     .././grub-core/lib/gnulib/memrchr.c \
     .././grub-core/lib/gnulib/msvc-inval.c \
     .././grub-core/lib/gnulib/msvc-inval.h \
     .././grub-core/lib/gnulib/msvc-nothrow.c \
     .././grub-core/lib/gnulib/msvc-nothrow.h \
     .././grub-core/lib/gnulib/nl_langinfo.c \
     .././grub-core/lib/gnulib/open.c \
     .././grub-core/lib/gnulib/openat-die.c \
     .././grub-core/lib/gnulib/openat-priv.h \
     .././grub-core/lib/gnulib/openat-proc.c \
     .././grub-core/lib/gnulib/openat.c \
     .././grub-core/lib/gnulib/openat.h \
     .././grub-core/lib/gnulib/pathmax.h \
     .././grub-core/lib/gnulib/pipe-safer.c \
     .././grub-core/lib/gnulib/printf-args.c \
     .././grub-core/lib/gnulib/printf-args.h \
     .././grub-core/lib/gnulib/printf-parse.c \
     .././grub-core/lib/gnulib/printf-parse.h \
     .././grub-core/lib/gnulib/progname.c \
     .././grub-core/lib/gnulib/progname.h \
     .././grub-core/lib/gnulib/rawmemchr.c \
     .././grub-core/lib/gnulib/realloc.c \
     .././grub-core/lib/gnulib/regcomp.c \
     .././grub-core/lib/gnulib/regex.c \
     .././grub-core/lib/gnulib/regex.h \
     .././grub-core/lib/gnulib/regex_internal.c \
     .././grub-core/lib/gnulib/regex_internal.h \
     .././grub-core/lib/gnulib/regexec.c \
     .././grub-core/lib/gnulib/save-cwd.c \
     .././grub-core/lib/gnulib/save-cwd.h \
     .././grub-core/lib/gnulib/size_max.h \
     .././grub-core/lib/gnulib/sleep.c \
     .././grub-core/lib/gnulib/stat-time.c \
     .././grub-core/lib/gnulib/stat-time.h \
     .././grub-core/lib/gnulib/stat-w32.c \
     .././grub-core/lib/gnulib/stat-w32.h \
     .././grub-core/lib/gnulib/stat.c \
     .././grub-core/lib/gnulib/stdalign.in.h \
     .././grub-core/lib/gnulib/stdbool.in.h \
     .././grub-core/lib/gnulib/stddef.in.h \
     .././grub-core/lib/gnulib/stdint.in.h \
     .././grub-core/lib/gnulib/stdio.in.h \
     .././grub-core/lib/gnulib/stdlib.in.h \
     .././grub-core/lib/gnulib/strcasecmp.c \
     .././grub-core/lib/gnulib/strchrnul.c \
     .././grub-core/lib/gnulib/strdup.c \
     .././grub-core/lib/gnulib/streq.h \
     .././grub-core/lib/gnulib/strerror-override.c \
     .././grub-core/lib/gnulib/strerror-override.h \
     .././grub-core/lib/gnulib/strerror.c \
     .././grub-core/lib/gnulib/string.in.h \
     .././grub-core/lib/gnulib/strings.in.h \
     .././grub-core/lib/gnulib/stripslash.c \
     .././grub-core/lib/gnulib/strncasecmp.c \
     .././grub-core/lib/gnulib/strndup.c \
     .././grub-core/lib/gnulib/strnlen.c \
     .././grub-core/lib/gnulib/strnlen1.c \
     .././grub-core/lib/gnulib/strnlen1.h \
     .././grub-core/lib/gnulib/sys_stat.in.h \
     .././grub-core/lib/gnulib/sys_types.in.h \
     .././grub-core/lib/gnulib/sysexits.in.h \
     .././grub-core/lib/gnulib/time.in.h \
     .././grub-core/lib/gnulib/unistd--.h \
     .././grub-core/lib/gnulib/unistd-safer.h \
     .././grub-core/lib/gnulib/unistd.c \
     .././grub-core/lib/gnulib/unistd.in.h \
     .././grub-core/lib/gnulib/unitypes.in.h \
     .././grub-core/lib/gnulib/uniwidth.in.h \
     .././grub-core/lib/gnulib/uniwidth/cjk.h \
     .././grub-core/lib/gnulib/uniwidth/width.c \
     .././grub-core/lib/gnulib/vasnprintf.c \
     .././grub-core/lib/gnulib/vasnprintf.h \
     .././grub-core/lib/gnulib/verify.h \
     .././grub-core/lib/gnulib/vsnprintf.c \
     .././grub-core/lib/gnulib/warn-on-use.h \
     .././grub-core/lib/gnulib/wchar.in.h \
     .././grub-core/lib/gnulib/wcrtomb.c \
     .././grub-core/lib/gnulib/wctype-h.c \
     .././grub-core/lib/gnulib/wctype.in.h \
     .././grub-core/lib/gnulib/wcwidth.c \
     .././grub-core/lib/gnulib/xalloc-oversized.h \
     .././grub-core/lib/gnulib/xsize.c \
     .././grub-core/lib/gnulib/xsize.h \
     .././grub-core/lib/hexdump.c \
     .././grub-core/lib/i386/backtrace.c \
     .././grub-core/lib/i386/halt.c \
     .././grub-core/lib/i386/pc/biosnum.c \
     .././grub-core/lib/i386/pc/vesa_modes_table.c \
     .././grub-core/lib/i386/random.c \
     .././grub-core/lib/i386/reboot.c \
     .././grub-core/lib/i386/relocator.c \
     .././grub-core/lib/i386/relocator_common_c.c \
     .././grub-core/lib/ieee1275/cmos.c \
     .././grub-core/lib/ieee1275/datetime.c \
     .././grub-core/lib/ieee1275/halt.c \
     .././grub-core/lib/ieee1275/reboot.c \
     .././grub-core/lib/ieee1275/relocator.c \
     .././grub-core/lib/legacy_parse.c \
     .././grub-core/lib/libgcrypt/cipher/ac.c \
     .././grub-core/lib/libgcrypt/cipher/arcfour.c \
     .././grub-core/lib/libgcrypt/cipher/bithelp.h \
     .././grub-core/lib/libgcrypt/cipher/blowfish.c \
     .././grub-core/lib/libgcrypt/cipher/bufhelp.h \
     .././grub-core/lib/libgcrypt/cipher/camellia-glue.c \
     .././grub-core/lib/libgcrypt/cipher/camellia.c \
     .././grub-core/lib/libgcrypt/cipher/camellia.h \
     .././grub-core/lib/libgcrypt/cipher/cast5.c \
     .././grub-core/lib/libgcrypt/cipher/cipher.c \
     .././grub-core/lib/libgcrypt/cipher/crc.c \
     .././grub-core/lib/libgcrypt/cipher/des.c \
     .././grub-core/lib/libgcrypt/cipher/dsa.c \
     .././grub-core/lib/libgcrypt/cipher/ecc.c \
     .././grub-core/lib/libgcrypt/cipher/elgamal.c \
     .././grub-core/lib/libgcrypt/cipher/hash-common.c \
     .././grub-core/lib/libgcrypt/cipher/hash-common.h \
     .././grub-core/lib/libgcrypt/cipher/hmac-tests.c \
     .././grub-core/lib/libgcrypt/cipher/idea.c \
     .././grub-core/lib/libgcrypt/cipher/kdf.c \
     .././grub-core/lib/libgcrypt/cipher/md.c \
     .././grub-core/lib/libgcrypt/cipher/md4.c \
     .././grub-core/lib/libgcrypt/cipher/md5.c \
     .././grub-core/lib/libgcrypt/cipher/primegen.c \
     .././grub-core/lib/libgcrypt/cipher/pubkey.c \
     .././grub-core/lib/libgcrypt/cipher/rfc2268.c \
     .././grub-core/lib/libgcrypt/cipher/rijndael-tables.h \
     .././grub-core/lib/libgcrypt/cipher/rijndael.c \
     .././grub-core/lib/libgcrypt/cipher/rmd.h \
     .././grub-core/lib/libgcrypt/cipher/rmd160.c \
     .././grub-core/lib/libgcrypt/cipher/rsa.c \
     .././grub-core/lib/libgcrypt/cipher/seed.c \
     .././grub-core/lib/libgcrypt/cipher/serpent.c \
     .././grub-core/lib/libgcrypt/cipher/sha1.c \
     .././grub-core/lib/libgcrypt/cipher/sha256.c \
     .././grub-core/lib/libgcrypt/cipher/sha512.c \
     .././grub-core/lib/libgcrypt/cipher/test-getrusage.c \
     .././grub-core/lib/libgcrypt/cipher/tiger.c \
     .././grub-core/lib/libgcrypt/cipher/twofish.c \
     .././grub-core/lib/libgcrypt/cipher/whirlpool.c \
     .././grub-core/lib/libgcrypt/mpi/amd64/mpi-asm-defs.h \
     .././grub-core/lib/libgcrypt/mpi/ec.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpi-asm-defs.h \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-add1.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-lshift.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-mul1.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-mul2.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-mul3.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-rshift.c \
     .././grub-core/lib/libgcrypt/mpi/generic/mpih-sub1.c \
     .././grub-core/lib/libgcrypt/mpi/generic/udiv-w-sdiv.c \
     .././grub-core/lib/libgcrypt/mpi/i386/syntax.h \
     .././grub-core/lib/libgcrypt/mpi/longlong.h \
     .././grub-core/lib/libgcrypt/mpi/m68k/syntax.h \
     .././grub-core/lib/libgcrypt/mpi/mips3/mpi-asm-defs.h \
     .././grub-core/lib/libgcrypt/mpi/mpi-add.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-bit.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-cmp.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-div.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-gcd.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-inline.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-inline.h \
     .././grub-core/lib/libgcrypt/mpi/mpi-internal.h \
     .././grub-core/lib/libgcrypt/mpi/mpi-inv.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-mod.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-mpow.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-mul.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-pow.c \
     .././grub-core/lib/libgcrypt/mpi/mpi-scan.c \
     .././grub-core/lib/libgcrypt/mpi/mpicoder.c \
     .././grub-core/lib/libgcrypt/mpi/mpih-div.c \
     .././grub-core/lib/libgcrypt/mpi/mpih-mul.c \
     .././grub-core/lib/libgcrypt/mpi/mpiutil.c \
     .././grub-core/lib/libgcrypt/mpi/powerpc32/syntax.h \
     .././grub-core/lib/libgcrypt/src/ath.c \
     .././grub-core/lib/libgcrypt/src/ath.h \
     .././grub-core/lib/libgcrypt/src/cipher-proto.h \
     .././grub-core/lib/libgcrypt/src/cipher.h \
     .././grub-core/lib/libgcrypt/src/dumpsexp.c \
     .././grub-core/lib/libgcrypt/src/fips.c \
     .././grub-core/lib/libgcrypt/src/g10lib.h \
     .././grub-core/lib/libgcrypt/src/gcrypt-module.h \
     .././grub-core/lib/libgcrypt/src/gcryptrnd.c \
     .././grub-core/lib/libgcrypt/src/getrandom.c \
     .././grub-core/lib/libgcrypt/src/hmac256.c \
     .././grub-core/lib/libgcrypt/src/hmac256.h \
     .././grub-core/lib/libgcrypt/src/hwfeatures.c \
     .././grub-core/lib/libgcrypt/src/missing-string.c \
     .././grub-core/lib/libgcrypt/src/module.c \
     .././grub-core/lib/libgcrypt/src/mpi.h \
     .././grub-core/lib/libgcrypt/src/secmem.h \
     .././grub-core/lib/libgcrypt/src/sexp.c \
     .././grub-core/lib/libgcrypt/src/stdmem.c \
     .././grub-core/lib/libgcrypt/src/stdmem.h \
     .././grub-core/lib/libgcrypt/src/types.h \
     .././grub-core/lib/libgcrypt/src/visibility.c \
     .././grub-core/lib/libgcrypt/src/visibility.h \
     .././grub-core/lib/libgcrypt_wrap/cipher_wrap.h \
     .././grub-core/lib/libgcrypt_wrap/mem.c \
     .././grub-core/lib/minilzo/lzoconf.h \
     .././grub-core/lib/minilzo/lzodefs.h \
     .././grub-core/lib/minilzo/minilzo.c \
     .././grub-core/lib/minilzo/minilzo.h \
     .././grub-core/lib/mips/arc/reboot.c \
     .././grub-core/lib/mips/loongson/reboot.c \
     .././grub-core/lib/mips/qemu_mips/reboot.c \
     .././grub-core/lib/mips/relocator.c \
     .././grub-core/lib/pbkdf2.c \
     .././grub-core/lib/posix_wrap/assert.h \
     .././grub-core/lib/posix_wrap/ctype.h \
     .././grub-core/lib/posix_wrap/errno.h \
     .././grub-core/lib/posix_wrap/inttypes.h \
     .././grub-core/lib/posix_wrap/langinfo.h \
     .././grub-core/lib/posix_wrap/limits.h \
     .././grub-core/lib/posix_wrap/localcharset.h \
     .././grub-core/lib/posix_wrap/locale.h \
     .././grub-core/lib/posix_wrap/stdint.h \
     .././grub-core/lib/posix_wrap/stdio.h \
     .././grub-core/lib/posix_wrap/stdlib.h \
     .././grub-core/lib/posix_wrap/string.h \
     .././grub-core/lib/posix_wrap/sys/types.h \
     .././grub-core/lib/posix_wrap/unistd.h \
     .././grub-core/lib/posix_wrap/wchar.h \
     .././grub-core/lib/posix_wrap/wctype.h \
     .././grub-core/lib/powerpc/relocator.c \
     .././grub-core/lib/priority_queue.c \
     .././grub-core/lib/progress.c \
     .././grub-core/lib/random.c \
     .././grub-core/lib/reed_solomon.c \
     .././grub-core/lib/relocator.c \
     .././grub-core/lib/syslinux_parse.c \
     .././grub-core/lib/uboot/reboot.c \
     .././grub-core/lib/x86_64/efi/relocator.c \
     .././grub-core/lib/xen/datetime.c \
     .././grub-core/lib/xen/halt.c \
     .././grub-core/lib/xen/reboot.c \
     .././grub-core/lib/xen/relocator.c \
     .././grub-core/lib/xzembed/xz.h \
     .././grub-core/lib/xzembed/xz_config.h \
     .././grub-core/lib/xzembed/xz_dec_bcj.c \
     .././grub-core/lib/xzembed/xz_dec_lzma2.c \
     .././grub-core/lib/xzembed/xz_dec_stream.c \
     .././grub-core/lib/xzembed/xz_lzma2.h \
     .././grub-core/lib/xzembed/xz_private.h \
     .././grub-core/lib/xzembed/xz_stream.h \
     .././grub-core/lib/zstd/bitstream.h \
     .././grub-core/lib/zstd/compiler.h \
     .././grub-core/lib/zstd/cpu.h \
     .././grub-core/lib/zstd/debug.c \
     .././grub-core/lib/zstd/debug.h \
     .././grub-core/lib/zstd/entropy_common.c \
     .././grub-core/lib/zstd/error_private.c \
     .././grub-core/lib/zstd/error_private.h \
     .././grub-core/lib/zstd/fse.h \
     .././grub-core/lib/zstd/fse_decompress.c \
     .././grub-core/lib/zstd/huf.h \
     .././grub-core/lib/zstd/huf_decompress.c \
     .././grub-core/lib/zstd/mem.h \
     .././grub-core/lib/zstd/module.c \
     .././grub-core/lib/zstd/xxhash.c \
     .././grub-core/lib/zstd/xxhash.h \
     .././grub-core/lib/zstd/zstd.h \
     .././grub-core/lib/zstd/zstd_common.c \
     .././grub-core/lib/zstd/zstd_decompress.c \
     .././grub-core/lib/zstd/zstd_errors.h \
     .././grub-core/lib/zstd/zstd_internal.h \
     .././grub-core/loader/aout.c \
     .././grub-core/loader/arm/linux.c \
     .././grub-core/loader/arm64/linux.c \
     .././grub-core/loader/arm64/xen_boot.c \
     .././grub-core/loader/efi/appleloader.c \
     .././grub-core/loader/efi/chainloader.c \
     .././grub-core/loader/efi/fdt.c \
     .././grub-core/loader/i386/bsd.c \
     .././grub-core/loader/i386/bsd32.c \
     .././grub-core/loader/i386/bsd64.c \
     .././grub-core/loader/i386/bsdXX.c \
     .././grub-core/loader/i386/bsd_pagetable.c \
     .././grub-core/loader/i386/coreboot/chainloader.c \
     .././grub-core/loader/i386/linux.c \
     .././grub-core/loader/i386/multiboot_mbi.c \
     .././grub-core/loader/i386/pc/chainloader.c \
     .././grub-core/loader/i386/pc/freedos.c \
     .././grub-core/loader/i386/pc/linux.c \
     .././grub-core/loader/i386/pc/ntldr.c \
     .././grub-core/loader/i386/pc/plan9.c \
     .././grub-core/loader/i386/pc/pxechainloader.c \
     .././grub-core/loader/i386/pc/truecrypt.c \
     .././grub-core/loader/i386/xen.c \
     .././grub-core/loader/i386/xen_file.c \
     .././grub-core/loader/i386/xen_file32.c \
     .././grub-core/loader/i386/xen_file64.c \
     .././grub-core/loader/i386/xen_fileXX.c \
     .././grub-core/loader/i386/xnu.c \
     .././grub-core/loader/ia64/efi/linux.c \
     .././grub-core/loader/linux.c \
     .././grub-core/loader/lzss.c \
     .././grub-core/loader/macho.c \
     .././grub-core/loader/macho32.c \
     .././grub-core/loader/macho64.c \
     .././grub-core/loader/machoXX.c \
     .././grub-core/loader/mips/linux.c \
     .././grub-core/loader/multiboot.c \
     .././grub-core/loader/multiboot_elfxx.c \
     .././grub-core/loader/multiboot_mbi2.c \
     .././grub-core/loader/powerpc/ieee1275/linux.c \
     .././grub-core/loader/riscv/linux.c \
     .././grub-core/loader/sparc64/ieee1275/linux.c \
     .././grub-core/loader/xnu.c \
     .././grub-core/loader/xnu_resume.c \
     .././grub-core/mmap/efi/mmap.c \
     .././grub-core/mmap/i386/mmap.c \
     .././grub-core/mmap/i386/pc/mmap.c \
     .././grub-core/mmap/i386/uppermem.c \
     .././grub-core/mmap/mips/uppermem.c \
     .././grub-core/mmap/mmap.c \
     .././grub-core/net/arp.c \
     .././grub-core/net/bootp.c \
     .././grub-core/net/dns.c \
     .././grub-core/net/drivers/efi/efinet.c \
     .././grub-core/net/drivers/emu/emunet.c \
     .././grub-core/net/drivers/i386/pc/pxe.c \
     .././grub-core/net/drivers/ieee1275/ofnet.c \
     .././grub-core/net/drivers/uboot/ubootnet.c \
     .././grub-core/net/ethernet.c \
     .././grub-core/net/http.c \
     .././grub-core/net/icmp.c \
     .././grub-core/net/icmp6.c \
     .././grub-core/net/ip.c \
     .././grub-core/net/net.c \
     .././grub-core/net/netbuff.c \
     .././grub-core/net/tcp.c \
     .././grub-core/net/tftp.c \
     .././grub-core/net/udp.c \
     .././grub-core/normal/auth.c \
     .././grub-core/normal/autofs.c \
     .././grub-core/normal/charset.c \
     .././grub-core/normal/cmdline.c \
     .././grub-core/normal/color.c \
     .././grub-core/normal/completion.c \
     .././grub-core/normal/context.c \
     .././grub-core/normal/crypto.c \
     .././grub-core/normal/datetime.c \
     .././grub-core/normal/dyncmd.c \
     .././grub-core/normal/main.c \
     .././grub-core/normal/menu.c \
     .././grub-core/normal/menu_entry.c \
     .././grub-core/normal/menu_text.c \
     .././grub-core/normal/misc.c \
     .././grub-core/normal/term.c \
     .././grub-core/osdep/apple/getroot.c \
     .././grub-core/osdep/apple/hostdisk.c \
     .././grub-core/osdep/aros/config.c \
     .././grub-core/osdep/aros/getroot.c \
     .././grub-core/osdep/aros/hostdisk.c \
     .././grub-core/osdep/aros/relpath.c \
     .././grub-core/osdep/basic/compress.c \
     .././grub-core/osdep/basic/emunet.c \
     .././grub-core/osdep/basic/getroot.c \
     .././grub-core/osdep/basic/hostdisk.c \
     .././grub-core/osdep/basic/init.c \
     .././grub-core/osdep/basic/no_platform.c \
     .././grub-core/osdep/basic/ofpath.c \
     .././grub-core/osdep/basic/platform.c \
     .././grub-core/osdep/basic/random.c \
     .././grub-core/osdep/blocklist.c \
     .././grub-core/osdep/bsd/getroot.c \
     .././grub-core/osdep/bsd/hostdisk.c \
     .././grub-core/osdep/compress.c \
     .././grub-core/osdep/config.c \
     .././grub-core/osdep/cputime.c \
     .././grub-core/osdep/devmapper/getroot.c \
     .././grub-core/osdep/devmapper/hostdisk.c \
     .././grub-core/osdep/dl.c \
     .././grub-core/osdep/emuconsole.c \
     .././grub-core/osdep/emunet.c \
     .././grub-core/osdep/exec.c \
     .././grub-core/osdep/freebsd/getroot.c \
     .././grub-core/osdep/freebsd/hostdisk.c \
     .././grub-core/osdep/generic/blocklist.c \
     .././grub-core/osdep/getroot.c \
     .././grub-core/osdep/haiku/getroot.c \
     .././grub-core/osdep/haiku/hostdisk.c \
     .././grub-core/osdep/hostdisk.c \
     .././grub-core/osdep/hurd/getroot.c \
     .././grub-core/osdep/hurd/hostdisk.c \
     .././grub-core/osdep/init.c \
     .././grub-core/osdep/linux/blocklist.c \
     .././grub-core/osdep/linux/emunet.c \
     .././grub-core/osdep/linux/getroot.c \
     .././grub-core/osdep/linux/hostdisk.c \
     .././grub-core/osdep/linux/ofpath.c \
     .././grub-core/osdep/linux/platform.c \
     .././grub-core/osdep/ofpath.c \
     .././grub-core/osdep/password.c \
     .././grub-core/osdep/platform.c \
     .././grub-core/osdep/platform_unix.c \
     .././grub-core/osdep/random.c \
     .././grub-core/osdep/relpath.c \
     .././grub-core/osdep/sleep.c \
     .././grub-core/osdep/sun/getroot.c \
     .././grub-core/osdep/sun/hostdisk.c \
     .././grub-core/osdep/unix/compress.c \
     .././grub-core/osdep/unix/config.c \
     .././grub-core/osdep/unix/cputime.c \
     .././grub-core/osdep/unix/dl.c \
     .././grub-core/osdep/unix/emuconsole.c \
     .././grub-core/osdep/unix/exec.c \
     .././grub-core/osdep/unix/getroot.c \
     .././grub-core/osdep/unix/hostdisk.c \
     .././grub-core/osdep/unix/password.c \
     .././grub-core/osdep/unix/platform.c \
     .././grub-core/osdep/unix/random.c \
     .././grub-core/osdep/unix/relpath.c \
     .././grub-core/osdep/unix/sleep.c \
     .././grub-core/osdep/windows/blocklist.c \
     .././grub-core/osdep/windows/config.c \
     .././grub-core/osdep/windows/cputime.c \
     .././grub-core/osdep/windows/dl.c \
     .././grub-core/osdep/windows/emuconsole.c \
     .././grub-core/osdep/windows/getroot.c \
     .././grub-core/osdep/windows/hostdisk.c \
     .././grub-core/osdep/windows/init.c \
     .././grub-core/osdep/windows/password.c \
     .././grub-core/osdep/windows/platform.c \
     .././grub-core/osdep/windows/random.c \
     .././grub-core/osdep/windows/relpath.c \
     .././grub-core/osdep/windows/sleep.c \
     .././grub-core/partmap/acorn.c \
     .././grub-core/partmap/amiga.c \
     .././grub-core/partmap/apple.c \
     .././grub-core/partmap/bsdlabel.c \
     .././grub-core/partmap/dfly.c \
     .././grub-core/partmap/dvh.c \
     .././grub-core/partmap/gpt.c \
     .././grub-core/partmap/msdos.c \
     .././grub-core/partmap/plan.c \
     .././grub-core/partmap/sun.c \
     .././grub-core/partmap/sunpc.c \
     .././grub-core/parttool/msdospart.c \
     .././grub-core/script/argv.c \
     .././grub-core/script/execute.c \
     .././grub-core/script/function.c \
     .././grub-core/script/lexer.c \
     .././grub-core/script/main.c \
     .././grub-core/script/script.c \
     .././grub-core/term/arc/console.c \
     .././grub-core/term/arc/serial.c \
     .././grub-core/term/arm/cros.c \
     .././grub-core/term/arm/cros_ec.c \
     .././grub-core/term/arm/pl050.c \
     .././grub-core/term/at_keyboard.c \
     .././grub-core/term/efi/console.c \
     .././grub-core/term/efi/serial.c \
     .././grub-core/term/gfxterm.c \
     .././grub-core/term/gfxterm_background.c \
     .././grub-core/term/i386/coreboot/cbmemc.c \
     .././grub-core/term/i386/pc/console.c \
     .././grub-core/term/i386/pc/mda_text.c \
     .././grub-core/term/i386/pc/vga_text.c \
     .././grub-core/term/ieee1275/console.c \
     .././grub-core/term/ieee1275/escc.c \
     .././grub-core/term/ieee1275/serial.c \
     .././grub-core/term/morse.c \
     .././grub-core/term/ns8250.c \
     .././grub-core/term/ps2.c \
     .././grub-core/term/serial.c \
     .././grub-core/term/spkmodem.c \
     .././grub-core/term/terminfo.c \
     .././grub-core/term/tparm.c \
     .././grub-core/term/uboot/console.c \
     .././grub-core/term/usb_keyboard.c \
     .././grub-core/term/xen/console.c \
     .././grub-core/tests/bswap_test.c \
     .././grub-core/tests/checksums.h \
     .././grub-core/tests/cmdline_cat_test.c \
     .././grub-core/tests/cmp_test.c \
     .././grub-core/tests/ctz_test.c \
     .././grub-core/tests/div_test.c \
     .././grub-core/tests/example_functional_test.c \
     .././grub-core/tests/fake_input.c \
     .././grub-core/tests/gfxterm_menu.c \
     .././grub-core/tests/legacy_password_test.c \
     .././grub-core/tests/lib/functional_test.c \
     .././grub-core/tests/lib/test.c \
     .././grub-core/tests/mul_test.c \
     .././grub-core/tests/pbkdf2_test.c \
     .././grub-core/tests/setjmp_test.c \
     .././grub-core/tests/shift_test.c \
     .././grub-core/tests/signature_test.c \
     .././grub-core/tests/signatures.h \
     .././grub-core/tests/sleep_test.c \
     .././grub-core/tests/strtoull_test.c \
     .././grub-core/tests/test_blockarg.c \
     .././grub-core/tests/video_checksum.c \
     .././grub-core/tests/videotest_checksum.c \
     .././grub-core/tests/xnu_uuid_test.c \
     .././grub-core/video/bitmap.c \
     .././grub-core/video/bitmap_scale.c \
     .././grub-core/video/bochs.c \
     .././grub-core/video/capture.c \
     .././grub-core/video/cirrus.c \
     .././grub-core/video/colors.c \
     .././grub-core/video/coreboot/cbfb.c \
     .././grub-core/video/efi_gop.c \
     .././grub-core/video/efi_uga.c \
     .././grub-core/video/emu/sdl.c \
     .././grub-core/video/fb/fbblit.c \
     .././grub-core/video/fb/fbfill.c \
     .././grub-core/video/fb/fbutil.c \
     .././grub-core/video/fb/video_fb.c \
     .././grub-core/video/i386/pc/vbe.c \
     .././grub-core/video/i386/pc/vga.c \
     .././grub-core/video/ieee1275.c \
     .././grub-core/video/radeon_fuloong2e.c \
     .././grub-core/video/radeon_yeeloong3a.c \
     .././grub-core/video/readers/jpeg.c \
     .././grub-core/video/readers/png.c \
     .././grub-core/video/readers/tga.c \
     .././grub-core/video/sis315_init.c \
     .././grub-core/video/sis315pro.c \
     .././grub-core/video/sm712.c \
     .././grub-core/video/sm712_init.c \
     .././grub-core/video/video.c \
     .././include/grub/acorn_filecore.h \
     .././include/grub/acpi.h \
     .././include/grub/aout.h \
     .././include/grub/arc/arc.h \
     .././include/grub/arc/console.h \
     .././include/grub/archelp.h \
     .././include/grub/arm/coreboot/console.h \
     .././include/grub/arm/coreboot/kernel.h \
     .././include/grub/arm/cros_ec.h \
     .././include/grub/arm/efi/memory.h \
     .././include/grub/arm/linux.h \
     .././include/grub/arm/reloc.h \
     .././include/grub/arm/setjmp.h \
     .././include/grub/arm/startup.h \
     .././include/grub/arm/system.h \
     .././include/grub/arm/time.h \
     .././include/grub/arm/types.h \
     .././include/grub/arm/uboot/kernel.h \
     .././include/grub/arm64/efi/memory.h \
     .././include/grub/arm64/linux.h \
     .././include/grub/arm64/reloc.h \
     .././include/grub/arm64/setjmp.h \
     .././include/grub/arm64/time.h \
     .././include/grub/arm64/types.h \
     .././include/grub/at_keyboard.h \
     .././include/grub/ata.h \
     .././include/grub/auth.h \
     .././include/grub/autoefi.h \
     .././include/grub/backtrace.h \
     .././include/grub/bitmap.h \
     .././include/grub/bitmap_scale.h \
     .././include/grub/boottime.h \
     .././include/grub/bsdlabel.h \
     .././include/grub/btrfs.h \
     .././include/grub/bufio.h \
     .././include/grub/cache.h \
     .././include/grub/cbfs_core.h \
     .././include/grub/charset.h \
     .././include/grub/cmos.h \
     .././include/grub/color.h \
     .././include/grub/command.h \
     .././include/grub/compiler-rt-emu.h \
     .././include/grub/compiler-rt.h \
     .././include/grub/compiler.h \
     .././include/grub/coreboot/lbio.h \
     .././include/grub/crypto.h \
     .././include/grub/cryptodisk.h \
     .././include/grub/cs5536.h \
     .././include/grub/datetime.h \
     .././include/grub/decompressor.h \
     .././include/grub/deflate.h \
     .././include/grub/device.h \
     .././include/grub/disk.h \
     .././include/grub/diskfilter.h \
     .././include/grub/dl.h \
     .././include/grub/dma.h \
     .././include/grub/efi/api.h \
     .././include/grub/efi/console.h \
     .././include/grub/efi/console_control.h \
     .././include/grub/efi/disk.h \
     .././include/grub/efi/edid.h \
     .././include/grub/efi/efi.h \
     .././include/grub/efi/fdtload.h \
     .././include/grub/efi/graphics_output.h \
     .././include/grub/efi/memory.h \
     .././include/grub/efi/pci.h \
     .././include/grub/efi/pe32.h \
     .././include/grub/efi/tpm.h \
     .././include/grub/efi/uga_draw.h \
     .././include/grub/efiemu/efiemu.h \
     .././include/grub/efiemu/runtime.h \
     .././include/grub/elf.h \
     .././include/grub/elfload.h \
     .././include/grub/emu/config.h \
     .././include/grub/emu/console.h \
     .././include/grub/emu/exec.h \
     .././include/grub/emu/getroot.h \
     .././include/grub/emu/hostdisk.h \
     .././include/grub/emu/hostfile.h \
     .././include/grub/emu/misc.h \
     .././include/grub/emu/net.h \
     .././include/grub/env.h \
     .././include/grub/env_private.h \
     .././include/grub/err.h \
     .././include/grub/exfat.h \
     .././include/grub/extcmd.h \
     .././include/grub/fat.h \
     .././include/grub/fbblit.h \
     .././include/grub/fbfill.h \
     .././include/grub/fbutil.h \
     .././include/grub/fdt.h \
     .././include/grub/fdtbus.h \
     .././include/grub/file.h \
     .././include/grub/fileid.h \
     .././include/grub/font.h \
     .././include/grub/fontformat.h \
     .././include/grub/fs.h \
     .././include/grub/fshelp.h \
     .././include/grub/gcry/types.h \
     .././include/grub/gcrypt/gpg-error.h \
     .././include/grub/gdb.h \
     .././include/grub/gfxmenu_model.h \
     .././include/grub/gfxmenu_view.h \
     .././include/grub/gfxterm.h \
     .././include/grub/gfxwidgets.h \
     .././include/grub/gpt_partition.h \
     .././include/grub/gui.h \
     .././include/grub/gui_string_util.h \
     .././include/grub/hfs.h \
     .././include/grub/hfsplus.h \
     .././include/grub/i18n.h \
     .././include/grub/i386/at_keyboard.h \
     .././include/grub/i386/bsd.h \
     .././include/grub/i386/cmos.h \
     .././include/grub/i386/coreboot/boot.h \
     .././include/grub/i386/coreboot/console.h \
     .././include/grub/i386/coreboot/kernel.h \
     .././include/grub/i386/coreboot/memory.h \
     .././include/grub/i386/coreboot/serial.h \
     .././include/grub/i386/coreboot/time.h \
     .././include/grub/i386/cpuid.h \
     .././include/grub/i386/efi/kernel.h \
     .././include/grub/i386/efi/memory.h \
     .././include/grub/i386/efi/serial.h \
     .././include/grub/i386/efiemu.h \
     .././include/grub/i386/floppy.h \
     .././include/grub/i386/freebsd_linker.h \
     .././include/grub/i386/freebsd_reboot.h \
     .././include/grub/i386/gdb.h \
     .././include/grub/i386/ieee1275/ieee1275.h \
     .././include/grub/i386/ieee1275/kernel.h \
     .././include/grub/i386/ieee1275/memory.h \
     .././include/grub/i386/ieee1275/serial.h \
     .././include/grub/i386/io.h \
     .././include/grub/i386/linux.h \
     .././include/grub/i386/macho.h \
     .././include/grub/i386/memory.h \
     .././include/grub/i386/memory_raw.h \
     .././include/grub/i386/multiboot.h \
     .././include/grub/i386/multiboot/boot.h \
     .././include/grub/i386/multiboot/console.h \
     .././include/grub/i386/multiboot/kernel.h \
     .././include/grub/i386/multiboot/memory.h \
     .././include/grub/i386/multiboot/serial.h \
     .././include/grub/i386/multiboot/time.h \
     .././include/grub/i386/netbsd_bootinfo.h \
     .././include/grub/i386/netbsd_reboot.h \
     .././include/grub/i386/openbsd_bootarg.h \
     .././include/grub/i386/openbsd_reboot.h \
     .././include/grub/i386/pc/apm.h \
     .././include/grub/i386/pc/biosdisk.h \
     .././include/grub/i386/pc/biosnum.h \
     .././include/grub/i386/pc/boot.h \
     .././include/grub/i386/pc/chainloader.h \
     .././include/grub/i386/pc/console.h \
     .././include/grub/i386/pc/int.h \
     .././include/grub/i386/pc/int_types.h \
     .././include/grub/i386/pc/kernel.h \
     .././include/grub/i386/pc/memory.h \
     .././include/grub/i386/pc/pxe.h \
     .././include/grub/i386/pc/time.h \
     .././include/grub/i386/pc/vbe.h \
     .././include/grub/i386/pc/vesa_modes_table.h \
     .././include/grub/i386/pci.h \
     .././include/grub/i386/pit.h \
     .././include/grub/i386/pmtimer.h \
     .././include/grub/i386/qemu/boot.h \
     .././include/grub/i386/qemu/console.h \
     .././include/grub/i386/qemu/kernel.h \
     .././include/grub/i386/qemu/memory.h \
     .././include/grub/i386/qemu/serial.h \
     .././include/grub/i386/qemu/time.h \
     .././include/grub/i386/rdmsr.h \
     .././include/grub/i386/reboot.h \
     .././include/grub/i386/relocator.h \
     .././include/grub/i386/relocator_private.h \
     .././include/grub/i386/setjmp.h \
     .././include/grub/i386/time.h \
     .././include/grub/i386/tsc.h \
     .././include/grub/i386/types.h \
     .././include/grub/i386/wrmsr.h \
     .././include/grub/i386/xen/hypercall.h \
     .././include/grub/i386/xen/kernel.h \
     .././include/grub/i386/xen/memory.h \
     .././include/grub/i386/xen_pvh/boot.h \
     .././include/grub/i386/xen_pvh/console.h \
     .././include/grub/i386/xen_pvh/int.h \
     .././include/grub/i386/xen_pvh/kernel.h \
     .././include/grub/i386/xen_pvh/memory.h \
     .././include/grub/i386/xen_pvh/time.h \
     .././include/grub/i386/xnu.h \
     .././include/grub/ia64/efi/memory.h \
     .././include/grub/ia64/efi/time.h \
     .././include/grub/ia64/kernel.h \
     .././include/grub/ia64/reloc.h \
     .././include/grub/ia64/setjmp.h \
     .././include/grub/ia64/time.h \
     .././include/grub/ia64/types.h \
     .././include/grub/icon_manager.h \
     .././include/grub/ieee1275/console.h \
     .././include/grub/ieee1275/ieee1275.h \
     .././include/grub/ieee1275/obdisk.h \
     .././include/grub/ieee1275/ofdisk.h \
     .././include/grub/kernel.h \
     .././include/grub/keyboard_layouts.h \
     .././include/grub/legacy_parse.h \
     .././include/grub/lib/LzFind.h \
     .././include/grub/lib/LzHash.h \
     .././include/grub/lib/LzmaDec.h \
     .././include/grub/lib/LzmaEnc.h \
     .././include/grub/lib/LzmaTypes.h \
     .././include/grub/lib/arg.h \
     .././include/grub/lib/cmdline.h \
     .././include/grub/lib/crc.h \
     .././include/grub/lib/envblk.h \
     .././include/grub/lib/hexdump.h \
     .././include/grub/libpciaccess.h \
     .././include/grub/linux.h \
     .././include/grub/list.h \
     .././include/grub/loader.h \
     .././include/grub/lvm.h \
     .././include/grub/macho.h \
     .././include/grub/machoload.h \
     .././include/grub/memory.h \
     .././include/grub/menu.h \
     .././include/grub/menu_viewer.h \
     .././include/grub/mips/arc/kernel.h \
     .././include/grub/mips/arc/memory.h \
     .././include/grub/mips/arc/time.h \
     .././include/grub/mips/asm.h \
     .././include/grub/mips/at_keyboard.h \
     .././include/grub/mips/cmos.h \
     .././include/grub/mips/io.h \
     .././include/grub/mips/kernel.h \
     .././include/grub/mips/loongson.h \
     .././include/grub/mips/loongson/at_keyboard.h \
     .././include/grub/mips/loongson/cmos.h \
     .././include/grub/mips/loongson/ec.h \
     .././include/grub/mips/loongson/kernel.h \
     .././include/grub/mips/loongson/memory.h \
     .././include/grub/mips/loongson/pci.h \
     .././include/grub/mips/loongson/serial.h \
     .././include/grub/mips/loongson/time.h \
     .././include/grub/mips/memory.h \
     .././include/grub/mips/mips.h \
     .././include/grub/mips/multiboot.h \
     .././include/grub/mips/pci.h \
     .././include/grub/mips/qemu_mips/at_keyboard.h \
     .././include/grub/mips/qemu_mips/cmos.h \
     .././include/grub/mips/qemu_mips/console.h \
     .././include/grub/mips/qemu_mips/kernel.h \
     .././include/grub/mips/qemu_mips/loader.h \
     .././include/grub/mips/qemu_mips/memory.h \
     .././include/grub/mips/qemu_mips/serial.h \
     .././include/grub/mips/qemu_mips/time.h \
     .././include/grub/mips/relocator.h \
     .././include/grub/mips/setjmp.h \
     .././include/grub/mips/time.h \
     .././include/grub/mips/types.h \
     .././include/grub/misc.h \
     .././include/grub/mm.h \
     .././include/grub/mm_private.h \
     .././include/grub/module_verifier.h \
     .././include/grub/msdos_partition.h \
     .././include/grub/multiboot.h \
     .././include/grub/multiboot2.h \
     .././include/grub/multiboot_loader.h \
     .././include/grub/net.h \
     .././include/grub/net/arp.h \
     .././include/grub/net/ethernet.h \
     .././include/grub/net/ip.h \
     .././include/grub/net/netbuff.h \
     .././include/grub/net/tcp.h \
     .././include/grub/net/udp.h \
     .././include/grub/normal.h \
     .././include/grub/ns8250.h \
     .././include/grub/ntfs.h \
     .././include/grub/offsets.h \
     .././include/grub/osdep/hostfile.h \
     .././include/grub/osdep/hostfile_aros.h \
     .././include/grub/osdep/hostfile_unix.h \
     .././include/grub/osdep/hostfile_windows.h \
     .././include/grub/parser.h \
     .././include/grub/partition.h \
     .././include/grub/parttool.h \
     .././include/grub/pci.h \
     .././include/grub/pciutils.h \
     .././include/grub/powerpc/ieee1275/ieee1275.h \
     .././include/grub/powerpc/kernel.h \
     .././include/grub/powerpc/memory.h \
     .././include/grub/powerpc/relocator.h \
     .././include/grub/powerpc/setjmp.h \
     .././include/grub/powerpc/time.h \
     .././include/grub/powerpc/types.h \
     .././include/grub/priority_queue.h \
     .././include/grub/procfs.h \
     .././include/grub/ps2.h \
     .././include/grub/pubkey.h \
     .././include/grub/random.h \
     .././include/grub/reader.h \
     .././include/grub/reed_solomon.h \
     .././include/grub/relocator.h \
     .././include/grub/relocator_private.h \
     .././include/grub/riscv32/efi/memory.h \
     .././include/grub/riscv32/linux.h \
     .././include/grub/riscv32/setjmp.h \
     .././include/grub/riscv32/time.h \
     .././include/grub/riscv32/types.h \
     .././include/grub/riscv64/efi/memory.h \
     .././include/grub/riscv64/linux.h \
     .././include/grub/riscv64/setjmp.h \
     .././include/grub/riscv64/time.h \
     .././include/grub/riscv64/types.h \
     .././include/grub/script_sh.h \
     .././include/grub/scsi.h \
     .././include/grub/scsicmd.h \
     .././include/grub/sdl.h \
     .././include/grub/search.h \
     .././include/grub/serial.h \
     .././include/grub/setjmp.h \
     .././include/grub/smbus.h \
     .././include/grub/sparc64/ieee1275/boot.h \
     .././include/grub/sparc64/ieee1275/ieee1275.h \
     .././include/grub/sparc64/ieee1275/kernel.h \
     .././include/grub/sparc64/setjmp.h \
     .././include/grub/sparc64/time.h \
     .././include/grub/sparc64/types.h \
     .././include/grub/speaker.h \
     .././include/grub/symbol.h \
     .././include/grub/syslinux_parse.h \
     .././include/grub/term.h \
     .././include/grub/terminfo.h \
     .././include/grub/test.h \
     .././include/grub/time.h \
     .././include/grub/tparm.h \
     .././include/grub/tpm.h \
     .././include/grub/trig.h \
     .././include/grub/types.h \
     .././include/grub/uboot/api_public.h \
     .././include/grub/uboot/console.h \
     .././include/grub/uboot/disk.h \
     .././include/grub/uboot/image.h \
     .././include/grub/uboot/uboot.h \
     .././include/grub/udf.h \
     .././include/grub/unicode.h \
     .././include/grub/usb.h \
     .././include/grub/usbdesc.h \
     .././include/grub/usbserial.h \
     .././include/grub/usbtrans.h \
     .././include/grub/util/install.h \
     .././include/grub/util/libnvpair.h \
     .././include/grub/util/libzfs.h \
     .././include/grub/util/misc.h \
     .././include/grub/util/mkimage.h \
     .././include/grub/util/ofpath.h \
     .././include/grub/util/resolve.h \
     .././include/grub/util/windows.h \
     .././include/grub/verify.h \
     .././include/grub/vga.h \
     .././include/grub/vgaregs.h \
     .././include/grub/video.h \
     .././include/grub/video_fb.h \
     .././include/grub/x86_64/at_keyboard.h \
     .././include/grub/x86_64/efi/boot.h \
     .././include/grub/x86_64/efi/kernel.h \
     .././include/grub/x86_64/efi/loader.h \
     .././include/grub/x86_64/efi/memory.h \
     .././include/grub/x86_64/efi/serial.h \
     .././include/grub/x86_64/io.h \
     .././include/grub/x86_64/linux.h \
     .././include/grub/x86_64/macho.h \
     .././include/grub/x86_64/memory.h \
     .././include/grub/x86_64/multiboot.h \
     .././include/grub/x86_64/pci.h \
     .././include/grub/x86_64/relocator.h \
     .././include/grub/x86_64/setjmp.h \
     .././include/grub/x86_64/time.h \
     .././include/grub/x86_64/types.h \
     .././include/grub/x86_64/xen/hypercall.h \
     .././include/grub/x86_64/xnu.h \
     .././include/grub/xen.h \
     .././include/grub/xen/relocator.h \
     .././include/grub/xen_file.h \
     .././include/grub/xnu.h \
     .././include/grub/zfs/dmu.h \
     .././include/grub/zfs/dmu_objset.h \
     .././include/grub/zfs/dnode.h \
     .././include/grub/zfs/dsl_dataset.h \
     .././include/grub/zfs/dsl_dir.h \
     .././include/grub/zfs/sa_impl.h \
     .././include/grub/zfs/spa.h \
     .././include/grub/zfs/uberblock_impl.h \
     .././include/grub/zfs/vdev_impl.h \
     .././include/grub/zfs/zap_impl.h \
     .././include/grub/zfs/zap_leaf.h \
     .././include/grub/zfs/zfs.h \
     .././include/grub/zfs/zfs_acl.h \
     .././include/grub/zfs/zfs_znode.h \
     .././include/grub/zfs/zil.h \
     .././include/grub/zfs/zio.h \
     .././include/grub/zfs/zio_checksum.h \
     .././include/multiboot.h \
     .././include/multiboot2.h \
     .././include/xen/arch-x86/xen-x86_32.h \
     .././include/xen/arch-x86/xen-x86_64.h \
     .././include/xen/arch-x86/xen.h \
     .././include/xen/elfnote.h \
     .././include/xen/event_channel.h \
     .././include/xen/grant_table.h \
     .././include/xen/hvm/hvm_op.h \
     .././include/xen/hvm/params.h \
     .././include/xen/hvm/start_info.h \
     .././include/xen/io/blkif.h \
     .././include/xen/io/console.h \
     .././include/xen/io/protocols.h \
     .././include/xen/io/ring.h \
     .././include/xen/io/xenbus.h \
     .././include/xen/io/xs_wire.h \
     .././include/xen/memory.h \
     .././include/xen/physdev.h \
     .././include/xen/sched.h \
     .././include/xen/trace.h \
     .././include/xen/xen-compat.h \
     .././include/xen/xen.h \
     .././tests/cmp_unit_test.c \
     .././tests/date_unit_test.c \
     .././tests/example_unit_test.c \
     .././tests/lib/unit_test.c \
     .././tests/printf_unit_test.c \
     .././util/bin2h.c \
     .././util/config.c \
     .././util/editenv.c \
     .././util/garbage-gen.c \
     .././util/getroot.c \
     .././util/glue-efi.c \
     .././util/grub-editenv.c \
     .././util/grub-file.c \
     .././util/grub-fstest.c \
     .././util/grub-glue-efi.c \
     .././util/grub-install-common.c \
     .././util/grub-install.c \
     .././util/grub-macbless.c \
     .././util/grub-macho2img.c \
     .././util/grub-menulst2cfg.c \
     .././util/grub-mkfont.c \
     .././util/grub-mkimage.c \
     .././util/grub-mkimage32.c \
     .././util/grub-mkimage64.c \
     .././util/grub-mkimagexx.c \
     .././util/grub-mklayout.c \
     .././util/grub-mknetdir.c \
     .././util/grub-mkpasswd-pbkdf2.c \
     .././util/grub-mkrelpath.c \
     .././util/grub-mkrescue.c \
     .././util/grub-mkstandalone.c \
     .././util/grub-module-verifier.c \
     .././util/grub-module-verifier32.c \
     .././util/grub-module-verifier64.c \
     .././util/grub-module-verifierXX.c \
     .././util/grub-mount.c \
     .././util/grub-pe2elf.c \
     .././util/grub-probe.c \
     .././util/grub-render-label.c \
     .././util/grub-script-check.c \
     .././util/grub-setup.c \
     .././util/grub-syslinux2cfg.c \
     .././util/ieee1275/grub-ofpathname.c \
     .././util/misc.c \
     .././util/mkimage.c \
     .././util/probe.c \
     .././util/render-label.c \
     .././util/resolve.c \
     .././util/setup.c \
     .././util/setup_bios.c \
     .././util/setup_sparc.c \
     .././util/spkmodem-recv.c

CATALOGS =  ast.gmo ca.gmo da.gmo de.gmo de@hebrew.gmo de_CH.gmo en@arabic.gmo en@cyrillic.gmo en@greek.gmo en@hebrew.gmo en@piglatin.gmo en@quot.gmo eo.gmo es.gmo fi.gmo fr.gmo gl.gmo hr.gmo hu.gmo id.gmo it.gmo ja.gmo ko.gmo lt.gmo nb.gmo nl.gmo pa.gmo pl.gmo pt.gmo pt_BR.gmo ro.gmo ru.gmo sl.gmo sr.gmo sv.gmo tr.gmo uk.gmo vi.gmo zh_CN.gmo zh_TW.gmo

# Makevars gets inserted here. (Don't remove this line!)
# Makefile variables for PO directory in any package using GNU gettext.

# Usually the message domain is the same as the package name.
DOMAIN = $(PACKAGE)

# These two variables depend on the location of this directory.
subdir = po
top_builddir = ..

# These options get passed to xgettext.
XGETTEXT_OPTIONS = --keyword=_ --keyword=N_ \
          \
 --flag=_:1:pass-c-format\
 --flag=N_:1:pass-c-format\
 --flag=error:3:c-format --flag=error_at_line:5:c-format\
\
 --from-code=UTF-8\
 $${end_of_xgettext_options+}

# This is the copyright holder that gets inserted into the header of the
# $(DOMAIN).pot file.  Set this to the copyright holder of the surrounding
# package.  (Note that the msgstr strings, extracted from the package's
# sources, belong to the copyright holder of the package.)  Translators are
# expected to transfer the copyright for their translations to this person
# or entity, or to disclaim their copyright.  The empty string stands for
# the public domain; in this case the translators are expected to disclaim
# their copyright.
COPYRIGHT_HOLDER = Free Software Foundation, Inc.

# This is the email address or URL to which the translators shall report
# bugs in the untranslated strings:
# - Strings which are not entire sentences, see the maintainer guidelines
#   in the GNU gettext documentation, section 'Preparing Strings'.
# - Strings which use unclear terms or require additional context to be
#   understood.
# - Strings which make invalid assumptions about notation of date, time or
#   money.
# - Pluralisation problems.
# - Incorrect English spelling.
# - Incorrect formatting.
# It can be your email address, or a mailing list address where translators
# can write to without being subscribed, or the URL of a web page through
# which the translators can contact you.
MSGID_BUGS_ADDRESS = bug-grub@gnu.org

# This is the list of locale categories, beyond LC_MESSAGES, for which the
# message catalogs shall be used.  It is usually empty.
EXTRA_LOCALE_CATEGORIES = 

# This tells whether the $(DOMAIN).pot file contains messages with an 'msgctxt'
# context.  Possible values are "yes" and "no".  Set this to yes if the
# package uses functions taking also a message context, like pgettext(), or
# if in $(XGETTEXT_OPTIONS) you define keywords with a context argument.
USE_MSGCTXT = no

# These options get passed to msgmerge.
# Useful options are in particular:
#   --previous            to keep previous msgids of translated messages,
#   --quiet               to reduce the verbosity.
MSGMERGE_OPTIONS =

.SUFFIXES:
.SUFFIXES: .po .gmo .mo .sed .sin .nop .po-create .po-update

.po.mo:
	@echo "$(MSGFMT) -c -o $@ $<"; \
	$(MSGFMT) -c -o t-$@ $< && mv t-$@ $@

.po.gmo:
	@lang=`echo $* | sed -e 's,.*/,,'`; \
	test "$(srcdir)" = . && cdcmd="" || cdcmd="cd $(srcdir) && "; \
	echo "$${cdcmd}rm -f $${lang}.gmo && $(GMSGFMT) -c --statistics --verbose -o $${lang}.gmo $${lang}.po"; \
	cd $(srcdir) && rm -f $${lang}.gmo && $(GMSGFMT) -c --statistics --verbose -o t-$${lang}.gmo $${lang}.po && mv t-$${lang}.gmo $${lang}.gmo

.sin.sed:
	sed -e '/^#/d' $< > t-$@
	mv t-$@ $@


all: all-no

all-yes: stamp-po
all-no:

# Ensure that the gettext macros and this Makefile.in.in are in sync.
CHECK_MACRO_VERSION = \
	test "$(GETTEXT_MACRO_VERSION)" = "0.18" \
	  || { echo "*** error: gettext infrastructure mismatch: using a Makefile.in.in from gettext version $(GETTEXT_MACRO_VERSION) but the autoconf macros are from gettext version 0.18" 1>&2; \
	       exit 1; \
	     }

# $(srcdir)/$(DOMAIN).pot is only created when needed. When xgettext finds no
# internationalized messages, no $(srcdir)/$(DOMAIN).pot is created (because
# we don't want to bother translators with empty POT files). We assume that
# LINGUAS is empty in this case, i.e. $(POFILES) and $(GMOFILES) are empty.
# In this case, stamp-po is a nop (i.e. a phony target).

# stamp-po is a timestamp denoting the last time at which the CATALOGS have
# been loosely updated. Its purpose is that when a developer or translator
# checks out the package via CVS, and the $(DOMAIN).pot file is not in CVS,
# "make" will update the $(DOMAIN).pot and the $(CATALOGS), but subsequent
# invocations of "make" will do nothing. This timestamp would not be necessary
# if updating the $(CATALOGS) would always touch them; however, the rule for
# $(POFILES) has been designed to not touch files that don't need to be
# changed.
stamp-po: $(srcdir)/$(DOMAIN).pot
	@$(CHECK_MACRO_VERSION)
	test ! -f $(srcdir)/$(DOMAIN).pot || \
	  test -z "$(GMOFILES)" || $(MAKE) $(GMOFILES)
	@test ! -f $(srcdir)/$(DOMAIN).pot || { \
	  echo "touch stamp-po" && \
	  echo timestamp > stamp-poT && \
	  mv stamp-poT stamp-po; \
	}

# Note: Target 'all' must not depend on target '$(DOMAIN).pot-update',
# otherwise packages like GCC can not be built if only parts of the source
# have been downloaded.

# This target rebuilds $(DOMAIN).pot; it is an expensive operation.
# Note that $(DOMAIN).pot is not touched if it doesn't need to be changed.
# The determination of whether the package xyz is a GNU one is based on the
# heuristic whether some file in the top level directory mentions "GNU xyz".
# If GNU 'find' is available, we avoid grepping through monster files.
$(DOMAIN).pot-update: $(POTFILES) $(srcdir)/POTFILES.in remove-potcdate.sed
	if { if (LC_ALL=C find --version) 2>/dev/null | grep GNU >/dev/null; then \
	       LC_ALL=C find -L $(top_srcdir) -maxdepth 1 -type f -size -10000000c -exec grep 'GNU grub' /dev/null '{}' ';' 2>/dev/null; \
	     else \
	       LC_ALL=C grep 'GNU grub' $(top_srcdir)/* 2>/dev/null; \
	     fi; \
	   } | grep -v 'libtool:' >/dev/null; then \
	  package_gnu='GNU '; \
	else \
	  package_gnu=''; \
	fi; \
	if test -n '$(MSGID_BUGS_ADDRESS)' || test '$(PACKAGE_BUGREPORT)' = '@'PACKAGE_BUGREPORT'@'; then \
	  msgid_bugs_address='$(MSGID_BUGS_ADDRESS)'; \
	else \
	  msgid_bugs_address='$(PACKAGE_BUGREPORT)'; \
	fi; \
	case `$(XGETTEXT) --version | sed 1q | sed -e 's,^[^0-9]*,,'` in \
	  '' | 0.[0-9] | 0.[0-9].* | 0.1[0-5] | 0.1[0-5].* | 0.16 | 0.16.[0-1]*) \
	    $(XGETTEXT) --default-domain=$(DOMAIN) --directory=$(top_srcdir) \
	      --add-comments=TRANSLATORS: $(XGETTEXT_OPTIONS)  --flag=argp_error:2:c-format --flag=argp_failure:4:c-format --flag=error:3:c-format --flag=error_at_line:5:c-format \
	      --files-from=$(srcdir)/POTFILES.in \
	      --copyright-holder='$(COPYRIGHT_HOLDER)' \
	      --msgid-bugs-address="$$msgid_bugs_address" \
	    ;; \
	  *) \
	    $(XGETTEXT) --default-domain=$(DOMAIN) --directory=$(top_srcdir) \
	      --add-comments=TRANSLATORS: $(XGETTEXT_OPTIONS)  --flag=argp_error:2:c-format --flag=argp_failure:4:c-format --flag=error:3:c-format --flag=error_at_line:5:c-format \
	      --files-from=$(srcdir)/POTFILES.in \
	      --copyright-holder='$(COPYRIGHT_HOLDER)' \
	      --package-name="$${package_gnu}grub" \
	      --package-version='2.04' \
	      --msgid-bugs-address="$$msgid_bugs_address" \
	    ;; \
	esac
	test ! -f $(DOMAIN).po || { \
	  if test -f $(srcdir)/$(DOMAIN).pot; then \
	    sed -f remove-potcdate.sed < $(srcdir)/$(DOMAIN).pot > $(DOMAIN).1po && \
	    sed -f remove-potcdate.sed < $(DOMAIN).po > $(DOMAIN).2po && \
	    if cmp $(DOMAIN).1po $(DOMAIN).2po >/dev/null 2>&1; then \
	      rm -f $(DOMAIN).1po $(DOMAIN).2po $(DOMAIN).po; \
	    else \
	      rm -f $(DOMAIN).1po $(DOMAIN).2po $(srcdir)/$(DOMAIN).pot && \
	      mv $(DOMAIN).po $(srcdir)/$(DOMAIN).pot; \
	    fi; \
	  else \
	    mv $(DOMAIN).po $(srcdir)/$(DOMAIN).pot; \
	  fi; \
	}

# This rule has no dependencies: we don't need to update $(DOMAIN).pot at
# every "make" invocation, only create it when it is missing.
# Only "make $(DOMAIN).pot-update" or "make dist" will force an update.
$(srcdir)/$(DOMAIN).pot:
	$(MAKE) $(DOMAIN).pot-update

# This target rebuilds a PO file if $(DOMAIN).pot has changed.
# Note that a PO file is not touched if it doesn't need to be changed.
$(POFILES): $(srcdir)/$(DOMAIN).pot
	@lang=`echo $@ | sed -e 's,.*/,,' -e 's/\.po$$//'`; \
	if test -f "$(srcdir)/$${lang}.po"; then \
	  test "$(srcdir)" = . && cdcmd="" || cdcmd="cd $(srcdir) && "; \
	  echo "$${cdcmd}$(MSGMERGE_UPDATE) $(MSGMERGE_OPTIONS) --lang=$${lang} $${lang}.po $(DOMAIN).pot"; \
	  cd $(srcdir) \
	    && { case `$(MSGMERGE_UPDATE) --version | sed 1q | sed -e 's,^[^0-9]*,,'` in \
	           '' | 0.[0-9] | 0.[0-9].* | 0.1[0-7] | 0.1[0-7].*) \
	             $(MSGMERGE_UPDATE) $(MSGMERGE_OPTIONS) $${lang}.po $(DOMAIN).pot;; \
	           *) \
	             $(MSGMERGE_UPDATE) $(MSGMERGE_OPTIONS) --lang=$${lang} $${lang}.po $(DOMAIN).pot;; \
	         esac; \
	       }; \
	else \
	  $(MAKE) $${lang}.po-create; \
	fi


install: install-exec install-data
install-exec:
install-data: install-data-no
	if test "$(PACKAGE)" = "gettext-tools"; then \
	  $(mkdir_p) $(DESTDIR)$(gettextsrcdir); \
	  for file in $(DISTFILES.common) Makevars.template; do \
	    $(INSTALL_DATA) $(srcdir)/$$file \
			    $(DESTDIR)$(gettextsrcdir)/$$file; \
	  done; \
	  for file in Makevars; do \
	    rm -f $(DESTDIR)$(gettextsrcdir)/$$file; \
	  done; \
	else \
	  : ; \
	fi
install-data-no: all
install-data-yes: all
	@catalogs='$(CATALOGS)'; \
	for cat in $$catalogs; do \
	  cat=`basename $$cat`; \
	  lang=`echo $$cat | sed -e 's/\.gmo$$//'`; \
	  dir=$(localedir)/$$lang/LC_MESSAGES; \
	  $(mkdir_p) $(DESTDIR)$$dir; \
	  if test -r $$cat; then realcat=$$cat; else realcat=$(srcdir)/$$cat; fi; \
	  $(INSTALL_DATA) $$realcat $(DESTDIR)$$dir/$(DOMAIN).mo; \
	  echo "installing $$realcat as $(DESTDIR)$$dir/$(DOMAIN).mo"; \
	  for lc in '' $(EXTRA_LOCALE_CATEGORIES); do \
	    if test -n "$$lc"; then \
	      if (cd $(DESTDIR)$(localedir)/$$lang && LC_ALL=C ls -l -d $$lc 2>/dev/null) | grep ' -> ' >/dev/null; then \
	        link=`cd $(DESTDIR)$(localedir)/$$lang && LC_ALL=C ls -l -d $$lc | sed -e 's/^.* -> //'`; \
	        mv $(DESTDIR)$(localedir)/$$lang/$$lc $(DESTDIR)$(localedir)/$$lang/$$lc.old; \
	        mkdir $(DESTDIR)$(localedir)/$$lang/$$lc; \
	        (cd $(DESTDIR)$(localedir)/$$lang/$$lc.old && \
	         for file in *; do \
	           if test -f $$file; then \
	             ln -s ../$$link/$$file $(DESTDIR)$(localedir)/$$lang/$$lc/$$file; \
	           fi; \
	         done); \
	        rm -f $(DESTDIR)$(localedir)/$$lang/$$lc.old; \
	      else \
	        if test -d $(DESTDIR)$(localedir)/$$lang/$$lc; then \
	          :; \
	        else \
	          rm -f $(DESTDIR)$(localedir)/$$lang/$$lc; \
	          mkdir $(DESTDIR)$(localedir)/$$lang/$$lc; \
	        fi; \
	      fi; \
	      rm -f $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo; \
	      ln -s ../LC_MESSAGES/$(DOMAIN).mo $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo 2>/dev/null || \
	      ln $(DESTDIR)$(localedir)/$$lang/LC_MESSAGES/$(DOMAIN).mo $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo 2>/dev/null || \
	      cp -p $(DESTDIR)$(localedir)/$$lang/LC_MESSAGES/$(DOMAIN).mo $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo; \
	      echo "installing $$realcat link as $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo"; \
	    fi; \
	  done; \
	done

install-strip: install

installdirs: installdirs-exec installdirs-data
installdirs-exec:
installdirs-data: installdirs-data-no
	if test "$(PACKAGE)" = "gettext-tools"; then \
	  $(mkdir_p) $(DESTDIR)$(gettextsrcdir); \
	else \
	  : ; \
	fi
installdirs-data-no:
installdirs-data-yes:
	@catalogs='$(CATALOGS)'; \
	for cat in $$catalogs; do \
	  cat=`basename $$cat`; \
	  lang=`echo $$cat | sed -e 's/\.gmo$$//'`; \
	  dir=$(localedir)/$$lang/LC_MESSAGES; \
	  $(mkdir_p) $(DESTDIR)$$dir; \
	  for lc in '' $(EXTRA_LOCALE_CATEGORIES); do \
	    if test -n "$$lc"; then \
	      if (cd $(DESTDIR)$(localedir)/$$lang && LC_ALL=C ls -l -d $$lc 2>/dev/null) | grep ' -> ' >/dev/null; then \
	        link=`cd $(DESTDIR)$(localedir)/$$lang && LC_ALL=C ls -l -d $$lc | sed -e 's/^.* -> //'`; \
	        mv $(DESTDIR)$(localedir)/$$lang/$$lc $(DESTDIR)$(localedir)/$$lang/$$lc.old; \
	        mkdir $(DESTDIR)$(localedir)/$$lang/$$lc; \
	        (cd $(DESTDIR)$(localedir)/$$lang/$$lc.old && \
	         for file in *; do \
	           if test -f $$file; then \
	             ln -s ../$$link/$$file $(DESTDIR)$(localedir)/$$lang/$$lc/$$file; \
	           fi; \
	         done); \
	        rm -f $(DESTDIR)$(localedir)/$$lang/$$lc.old; \
	      else \
	        if test -d $(DESTDIR)$(localedir)/$$lang/$$lc; then \
	          :; \
	        else \
	          rm -f $(DESTDIR)$(localedir)/$$lang/$$lc; \
	          mkdir $(DESTDIR)$(localedir)/$$lang/$$lc; \
	        fi; \
	      fi; \
	    fi; \
	  done; \
	done

# Define this as empty until I found a useful application.
installcheck:

uninstall: uninstall-exec uninstall-data
uninstall-exec:
uninstall-data: uninstall-data-no
	if test "$(PACKAGE)" = "gettext-tools"; then \
	  for file in $(DISTFILES.common) Makevars.template; do \
	    rm -f $(DESTDIR)$(gettextsrcdir)/$$file; \
	  done; \
	else \
	  : ; \
	fi
uninstall-data-no:
uninstall-data-yes:
	catalogs='$(CATALOGS)'; \
	for cat in $$catalogs; do \
	  cat=`basename $$cat`; \
	  lang=`echo $$cat | sed -e 's/\.gmo$$//'`; \
	  for lc in LC_MESSAGES $(EXTRA_LOCALE_CATEGORIES); do \
	    rm -f $(DESTDIR)$(localedir)/$$lang/$$lc/$(DOMAIN).mo; \
	  done; \
	done

check: all

info dvi ps pdf html tags TAGS ctags CTAGS ID:

mostlyclean:
	rm -f remove-potcdate.sed
	rm -f stamp-poT
	rm -f core core.* $(DOMAIN).po $(DOMAIN).1po $(DOMAIN).2po *.new.po
	rm -fr *.o

clean: mostlyclean

distclean: clean
	rm -f Makefile Makefile.in POTFILES *.mo

maintainer-clean: distclean
	@echo "This command is intended for maintainers to use;"
	@echo "it deletes files that may require special tools to rebuild."
	rm -f stamp-po $(GMOFILES)

distdir = $(top_builddir)/$(PACKAGE)-$(VERSION)/$(subdir)
dist distdir:
	$(MAKE) update-po
	@$(MAKE) dist2
# This is a separate target because 'update-po' must be executed before.
dist2: stamp-po $(DISTFILES)
	dists="$(DISTFILES)"; \
	if test "$(PACKAGE)" = "gettext-tools"; then \
	  dists="$$dists Makevars.template"; \
	fi; \
	if test -f $(srcdir)/$(DOMAIN).pot; then \
	  dists="$$dists $(DOMAIN).pot stamp-po"; \
	fi; \
	if test -f $(srcdir)/ChangeLog; then \
	  dists="$$dists ChangeLog"; \
	fi; \
	for i in 0 1 2 3 4 5 6 7 8 9; do \
	  if test -f $(srcdir)/ChangeLog.$$i; then \
	    dists="$$dists ChangeLog.$$i"; \
	  fi; \
	done; \
	if test -f $(srcdir)/LINGUAS; then dists="$$dists LINGUAS"; fi; \
	for file in $$dists; do \
	  if test -f $$file; then \
	    cp -p $$file $(distdir) || exit 1; \
	  else \
	    cp -p $(srcdir)/$$file $(distdir) || exit 1; \
	  fi; \
	done

update-po: Makefile
	$(MAKE) $(DOMAIN).pot-update
	test -z "$(UPDATEPOFILES)" || $(MAKE) $(UPDATEPOFILES)
	$(MAKE) update-gmo

# General rule for creating PO files.

.nop.po-create:
	@lang=`echo $@ | sed -e 's/\.po-create$$//'`; \
	echo "File $$lang.po does not exist. If you are a translator, you can create it through 'msginit'." 1>&2; \
	exit 1

# General rule for updating PO files.

.nop.po-update:
	@lang=`echo $@ | sed -e 's/\.po-update$$//'`; \
	if test "$(PACKAGE)" = "gettext-tools"; then PATH=`pwd`/../src:$$PATH; fi; \
	tmpdir=`pwd`; \
	echo "$$lang:"; \
	test "$(srcdir)" = . && cdcmd="" || cdcmd="cd $(srcdir) && "; \
	echo "$${cdcmd}$(MSGMERGE) $(MSGMERGE_OPTIONS) --lang=$$lang $$lang.po $(DOMAIN).pot -o $$lang.new.po"; \
	cd $(srcdir); \
	if { case `$(MSGMERGE) --version | sed 1q | sed -e 's,^[^0-9]*,,'` in \
	       '' | 0.[0-9] | 0.[0-9].* | 0.1[0-7] | 0.1[0-7].*) \
	         $(MSGMERGE) $(MSGMERGE_OPTIONS) -o $$tmpdir/$$lang.new.po $$lang.po $(DOMAIN).pot;; \
	       *) \
	         $(MSGMERGE) $(MSGMERGE_OPTIONS) --lang=$$lang -o $$tmpdir/$$lang.new.po $$lang.po $(DOMAIN).pot;; \
	     esac; \
	   }; then \
	  if cmp $$lang.po $$tmpdir/$$lang.new.po >/dev/null 2>&1; then \
	    rm -f $$tmpdir/$$lang.new.po; \
	  else \
	    if mv -f $$tmpdir/$$lang.new.po $$lang.po; then \
	      :; \
	    else \
	      echo "msgmerge for $$lang.po failed: cannot move $$tmpdir/$$lang.new.po to $$lang.po" 1>&2; \
	      exit 1; \
	    fi; \
	  fi; \
	else \
	  echo "msgmerge for $$lang.po failed!" 1>&2; \
	  rm -f $$tmpdir/$$lang.new.po; \
	fi

$(DUMMYPOFILES):

update-gmo: Makefile $(GMOFILES)
	@:

# Recreate Makefile by invoking config.status. Explicitly invoke the shell,
# because execution permission bits may not work on the current file system.
# Use /bin/sh, which is the shell determined by autoconf for the use by its
# scripts, not $(SHELL) which is hardwired to /bin/sh and may be deficient.
Makefile: Makefile.in.in Makevars $(top_builddir)/config.status POTFILES.in LINGUAS
	cd $(top_builddir) \
	  && /bin/sh ./config.status $(subdir)/$@.in po-directories

force:

# Tell versions [3.59,3.63) of GNU make not to export all variables.
# Otherwise a system limit (for SysV at least) may be exceeded.
.NOEXPORT:
# Special Makefile rules for Pig Latin message catalogs with quotation marks.

DISTFILES.common.extra3 = piglatin.sed en@piglatin.header insert-header.sin Rules-piglatin

.SUFFIXES: .insert-header .po-update-en-piglatin

en@piglatin.po-create:
	$(MAKE) en@piglatin.po-update

en@piglatin.po-update: en@piglatin.po-update-en

.insert-header.po-update-en-piglatin:
	@lang=`echo $@ | sed -e 's/\.po-update-en-piglatin$$//'`; \
	if test "$(PACKAGE)" = "gettext"; then PATH=`pwd`/../src:$$PATH; GETTEXTLIBDIR=`cd $(top_srcdir)/src && pwd`; export GETTEXTLIBDIR; fi; \
	tmpdir=`pwd`; \
	echo "$$lang:"; \
	ll=`echo $$lang | sed -e 's/@.*//'`; \
	LC_ALL=C; export LC_ALL; \
	cd $(srcdir); \
	if $(MSGINIT) -i $(DOMAIN).pot --no-translator -l $$lang -o - 2>/dev/null | sed -f $$tmpdir/$$lang.insert-header | $(MSGCONV) -t UTF-8 | $(MSGFILTER) sed -f `echo $$lang | sed -e 's/.*@//'`.sed 2>/dev/null > $$tmpdir/$$lang.new.po; then \
	  if cmp $$lang.po $$tmpdir/$$lang.new.po >/dev/null 2>&1; then \
	    rm -f $$tmpdir/$$lang.new.po; \
	  else \
	    if mv -f $$tmpdir/$$lang.new.po $$lang.po; then \
	      :; \
	    else \
	      echo "creation of $$lang.po failed: cannot move $$tmpdir/$$lang.new.po to $$lang.po" 1>&2; \
	      exit 1; \
	    fi; \
	  fi; \
	else \
	  echo "creation of $$lang.po failed!" 1>&2; \
	  rm -f $$tmpdir/$$lang.new.po; \
	fi

en@piglatin.insert-header: insert-header.sin
	sed -e '/^#/d' -e 's/HEADER/en@piglatin.header/g' $(srcdir)/insert-header.sin > en@piglatin.insert-header

mostlyclean: mostlyclean-piglatin
mostlyclean-piglatin:
	rm -f *.insert-header
# Special Makefile rules for English message catalogs with quotation marks.

DISTFILES.common.extra1 = quot.sed boldquot.sed en@quot.header en@boldquot.header insert-header.sin Rules-quot

.SUFFIXES: .insert-header .po-update-en

en@quot.po-create:
	$(MAKE) en@quot.po-update
en@boldquot.po-create:
	$(MAKE) en@boldquot.po-update

en@quot.po-update: en@quot.po-update-en
en@boldquot.po-update: en@boldquot.po-update-en

.insert-header.po-update-en:
	@lang=`echo $@ | sed -e 's/\.po-update-en$$//'`; \
	if test "$(PACKAGE)" = "gettext-tools"; then PATH=`pwd`/../src:$$PATH; GETTEXTLIBDIR=`cd $(top_srcdir)/src && pwd`; export GETTEXTLIBDIR; fi; \
	tmpdir=`pwd`; \
	echo "$$lang:"; \
	ll=`echo $$lang | sed -e 's/@.*//'`; \
	LC_ALL=C; export LC_ALL; \
	cd $(srcdir); \
	if $(MSGINIT) -i $(DOMAIN).pot --no-translator -l $$lang -o - 2>/dev/null | sed -f $$tmpdir/$$lang.insert-header | $(MSGCONV) -t UTF-8 | $(MSGFILTER) $(SED) -f `echo $$lang | sed -e 's/.*@//'`.sed 2>/dev/null > $$tmpdir/$$lang.new.po; then \
	  if cmp $$lang.po $$tmpdir/$$lang.new.po >/dev/null 2>&1; then \
	    rm -f $$tmpdir/$$lang.new.po; \
	  else \
	    if mv -f $$tmpdir/$$lang.new.po $$lang.po; then \
	      :; \
	    else \
	      echo "creation of $$lang.po failed: cannot move $$tmpdir/$$lang.new.po to $$lang.po" 1>&2; \
	      exit 1; \
	    fi; \
	  fi; \
	else \
	  echo "creation of $$lang.po failed!" 1>&2; \
	  rm -f $$tmpdir/$$lang.new.po; \
	fi

en@quot.insert-header: insert-header.sin
	sed -e '/^#/d' -e 's/HEADER/en@quot.header/g' $(srcdir)/insert-header.sin > en@quot.insert-header

en@boldquot.insert-header: insert-header.sin
	sed -e '/^#/d' -e 's/HEADER/en@boldquot.header/g' $(srcdir)/insert-header.sin > en@boldquot.insert-header

mostlyclean: mostlyclean-quot
mostlyclean-quot:
	rm -f *.insert-header
# Special Makefile rules for transforming German typography into Swiss one.

DISTFILES.common.extra2 = swiss.sed Rules-swiss

de_CH.po-create: de.po swiss.sed
	LC_ALL=C.UTF-8 $(MSGFILTER) -i $< -o $(srcdir)/de_CH.po sed -f $(srcdir)/swiss.sed

