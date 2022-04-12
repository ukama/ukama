gfxinit-y += hw-gfx-dp_aux_ch.adb
gfxinit-y += hw-gfx-dp_aux_ch.ads
gfxinit-y += hw-gfx-dp_defs.ads
gfxinit-y += hw-gfx-dp_info.adb
gfxinit-y += hw-gfx-dp_info.ads
gfxinit-y += hw-gfx-dp_training.adb
gfxinit-y += hw-gfx-dp_training.ads
gfxinit-y += hw-gfx-edid.adb
gfxinit-y += hw-gfx-edid.ads
gfxinit-y += hw-gfx-gma-config_helpers.adb
gfxinit-y += hw-gfx-gma-config_helpers.ads
gfxinit-y += hw-gfx-gma-connector_info.adb
gfxinit-y += hw-gfx-gma-connector_info.ads
gfxinit-y += hw-gfx-gma-connectors.ads
gfxinit-y += hw-gfx-gma-display_probing.adb
gfxinit-y += hw-gfx-gma-display_probing.ads
gfxinit-y += hw-gfx-gma-dp_aux_ch.ads
gfxinit-y += hw-gfx-gma-dp_aux_request.adb
gfxinit-y += hw-gfx-gma-dp_aux_request.ads
gfxinit-y += hw-gfx-gma-dp_info.ads
gfxinit-y += hw-gfx-gma-i2c.adb
gfxinit-y += hw-gfx-gma-i2c.ads
gfxinit-y += hw-gfx-gma-panel.adb
gfxinit-y += hw-gfx-gma-panel.ads
gfxinit-y += hw-gfx-gma-pch-fdi.adb
gfxinit-y += hw-gfx-gma-pch-fdi.ads
gfxinit-y += hw-gfx-gma-pch-sideband.adb
gfxinit-y += hw-gfx-gma-pch-sideband.ads
gfxinit-y += hw-gfx-gma-pch-transcoder.adb
gfxinit-y += hw-gfx-gma-pch-transcoder.ads
gfxinit-y += hw-gfx-gma-pch-vga.adb
gfxinit-y += hw-gfx-gma-pch-vga.ads
gfxinit-y += hw-gfx-gma-pch.ads
gfxinit-y += hw-gfx-gma-pcode.adb
gfxinit-y += hw-gfx-gma-pcode.ads
gfxinit-y += hw-gfx-gma-pipe_setup.adb
gfxinit-y += hw-gfx-gma-pipe_setup.ads
gfxinit-y += hw-gfx-gma-port_detect.ads
gfxinit-y += hw-gfx-gma-registers.adb
gfxinit-y += hw-gfx-gma-registers.ads
gfxinit-y += hw-gfx-gma-transcoder.adb
gfxinit-y += hw-gfx-gma-transcoder.ads
gfxinit-y += hw-gfx-gma.adb
gfxinit-y += hw-gfx-gma.ads
gfxinit-y += hw-gfx-i2c.ads
gfxinit-y += hw-gfx.ads
gfxinit-y += hw-gfx-framebuffer_filler.adb
gfxinit-y += hw-gfx-framebuffer_filler.ads

CONFIG_GFX_GMA_GENERATION	:= $(call strip_quotes,$(CONFIG_GFX_GMA_GENERATION))
CONFIG_GFX_GMA_CPU		:= $(call strip_quotes,$(CONFIG_GFX_GMA_CPU))
CONFIG_GFX_GMA_CPU_VARIANT	:= $(call strip_quotes,$(CONFIG_GFX_GMA_CPU_VARIANT))
CONFIG_GFX_GMA_INTERNAL_PORT	:= $(call strip_quotes,$(CONFIG_GFX_GMA_INTERNAL_PORT))
CONFIG_GFX_GMA_ANALOG_I2C_PORT	:= $(call strip_quotes,$(CONFIG_GFX_GMA_ANALOG_I2C_PORT))

_GEN_NONCONST := $(strip \
		 $(if $(filter G45,$(CONFIG_GFX_GMA_GENERATION)),g45, \
		 $(if $(filter Ironlake,$(CONFIG_GFX_GMA_GENERATION)),ilk, \
		 $(if $(filter Haswell,$(CONFIG_GFX_GMA_GENERATION)),hsw,  \
		 $(if $(filter Skylake,$(CONFIG_GFX_GMA_GENERATION)),skl)))))
# GNATprove (GPL 2017) doesn't realize when a boolean expression
# that depends both on static values and variables can be evalu-
# ated at compile time (e.g. `False and then Variable` is always
# `False` and GNAT acts appropriately). So for now, we generate
# functions instead of constant expressions for these mixed ex-
# pressions.
_GEN_CONST_TARGET := <cpufunc> # set to `constant` to generate constants.

hw-gfx-gma-config-ads := $(subst //,/,$(call src-to-obj,,$(dir)/hw-gfx-gma-config).ads)

ifeq ($(CONFIG_GFX_GMA_DYN_CPU),y)
$(hw-gfx-gma-config-ads): $(dir)/hw-gfx-gma-config.ads.template $(cnf)
	printf "    GENERATE   $(patsubst /%,%,$(subst $(obj)/,,$@))\n"
	sed -e's/<<GEN>>/$(CONFIG_GFX_GMA_GENERATION)/' \
	    -e's/<<INTERNAL_PORT>>/$(CONFIG_GFX_GMA_INTERNAL_PORT)/' \
	    -e's/<<ANALOG_I2C_PORT>>/$(CONFIG_GFX_GMA_ANALOG_I2C_PORT)/' \
	    -e's/<<DEFAULT_MMIO_BASE>>/$(CONFIG_GFX_GMA_DEFAULT_MMIO)/' \
	    -e'/constant Gen_CPU\(_Var\)\?/d' \
	    -e's/<genbool>/constant Boolean/' \
	    -e's/<\(\(g45\|ilk\|hsw\|skl\)\(...\)\?\)bool>/<\1var> Boolean/' \
	    $(if $(_GEN_NONCONST),-e's/<\(...\)\?$(_GEN_NONCONST)\(...\)\?var>/<cpufunc>/') \
	    -e's/<\(g45\|ilk\|hsw\|skl\)\(...\)\?var>/$(_GEN_CONST_TARGET)/' \
	    -e's/\(.*: *<cpufunc>.*:=\) *\(.*\);/\1\n     (\2);/' \
	    -e's/\([^ ]\+\) *: *<cpufunc> \+\([^ ]*\) *:=/function \1 return \2 is/' \
	    -e's/<cpunull>//' \
	    $< >$@
else
$(hw-gfx-gma-config-ads): $(dir)/hw-gfx-gma-config.ads.template $(cnf)
	printf "    GENERATE   $(patsubst /%,%,$(subst $(obj)/,,$@))\n"
	sed -e's/<<GEN>>/$(CONFIG_GFX_GMA_GENERATION)/' \
	    -e's/<<CPU>>/$(CONFIG_GFX_GMA_CPU)/' \
	    -e's/<<CPU_VARIANT>>/$(CONFIG_GFX_GMA_CPU_VARIANT)/' \
	    -e's/<<INTERNAL_PORT>>/$(CONFIG_GFX_GMA_INTERNAL_PORT)/' \
	    -e's/<<ANALOG_I2C_PORT>>/$(CONFIG_GFX_GMA_ANALOG_I2C_PORT)/' \
	    -e's/<<DEFAULT_MMIO_BASE>>/$(CONFIG_GFX_GMA_DEFAULT_MMIO)/' \
	    -e":s$$(printf '\n ')/,$$/{N;s/,\n.*Dyn_CPU\(_Var\)\?[^,)]*//;ts$$(printf '\n ')P;D;}" \
	    -e'/Dyn_CPU\(_Var\)\?/d' \
	    -e's/<\(gen\|\(g45\|ilk\|hsw\|skl\)\(...\)\?\)bool>/constant Boolean/' \
	    -e's/<\(\(g45\|ilk\|hsw\|skl\)\(...\)\?\)var>/constant/' \
	    -e's/<cpunull>/ is null/' \
	    $< >$@
endif
gfxinit-gen-y += $(hw-gfx-gma-config-ads)
gfxinit-$(CONFIG_GFX_GMA_DYN_CPU) += dyncpu/hw-gfx-gma-config.adb

ifneq ($(filter G45,$(CONFIG_GFX_GMA_GENERATION)),)
subdirs-y += g45
else ifneq ($(filter Ironlake,$(CONFIG_GFX_GMA_GENERATION)),)
subdirs-y += ironlake
else ifneq ($(filter Haswell,$(CONFIG_GFX_GMA_GENERATION)),)
subdirs-y += haswell_shared haswell
else ifneq ($(filter Broxton,$(CONFIG_GFX_GMA_GENERATION)),)
subdirs-y += haswell_shared broxton
else ifneq ($(filter Skylake,$(CONFIG_GFX_GMA_GENERATION)),)
subdirs-y += haswell_shared skylake
endif
