cmake_minimum_required(VERSION 2.8)

include(ExternalProject)
set(SYSROOT /Volumes/Skivavbild/boxeebox-xbmc/tools/boxeebox/sysroot)
set(BUILD_DEST ${CMAKE_BINARY_DIR}/target)
set(ENV{BUILD_DEST} ${BUILD_DEST})
set(ENV{KERNEL_VER} linux-2.6.35)
set(ENV{CROSS_COMPILE} ${BUILD_DEST}/bin/i686-cm-linux-)
set(ENV{ARCH} i386)

if(${VERBOSE})
	set(ENV{LOUD} 1)
endif()



add_custom_command(
	OUTPUT ${BUILD_DEST}/bin
	COMMAND ${CMAKE_COMMAND} -E make_directory ${BUILD_DEST}/bin
)

set(TARGETCMDS "")
set(CMDS addr2line ar as c++ cpp g++ gcc gccbug gcov gprof ld nm objcopy objdump ranlib readelf size strings strip)
foreach(CMD ${CMDS})
	string(TOUPPER TARGET${CMD} t)
	string(REPLACE ++ "XX" t ${t})
	set(ENV{${t}}  $ENV{CROSS_COMPILE}${CMD})
	set(TARGETCMDS ${TARGETCMDS} "${t}=$ENV{${t}}")
	#message("${t} = $ENV{${t}}")
	add_custom_target(
		${CMD}
		COMMAND ${CMAKE_COMMAND} -E create_symlink ${SYSROOT}/usr/bin/i686-pc-linux-gnu-${CMD} $ENV{${t}}
		DEPENDS ${BUILD_DEST}/bin
	)
endforeach()
set(ENV{TARGETCC} ${BUILD_DEST}/bin/i686-cm-linux-gcc)
set(TARGETCMDS ${TARGETCMDS} TARGETCC=$ENV{TARGETCC})
set(SMD_MAKE ARCH=$ENV{ARCH} CROSS_COMPILE=$ENV{CROSS_COMPILE} ${TARGETCMDS} BUILD_DEST=$ENV{BUILD_DEST} KERNEL_VER=$ENV{KERNEL_VER} make all)
#message("SMD_MAKE = ${SMD_MAKE}")

add_custom_target(
	toolchain
	DEPENDS ${CMDS}
)

ExternalProject_Add(
	smd_common
	URL 				${CMAKE_SOURCE_DIR}/SMD_Common.tar.gz
	SOURCE_DIR 			SMD_Common
	BINARY_DIR 			SMD_Common/SMD_Common
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS toolchain
)



ExternalProject_Add(
	kernel
	URL 				${CMAKE_SOURCE_DIR}/kernel.tar.gz
	SOURCE_DIR			kernel
	BINARY_DIR			kernel/kernel
	PATCH_COMMAND
		COMMAND ${CMAKE_COMMAND} -E make_directory kernel/linux-2.6.35/.git
		COMMAND ${CMAKE_COMMAND} -E create_symlink /usr/include/malloc/malloc.h kernel/linux-2.6.35/scripts/genksyms/malloc.h
		COMMAND ${CMAKE_COMMAND} -E copy ${SYSROOT}/usr/include/elf.h kernel/linux-2.6.35/scripts/mod/elf.h
		COMMAND ${CMAKE_COMMAND} -E copy ${SYSROOT}/usr/include/sys/sysmacros.h kernel/linux-2.6.35/arch/x86/boot/tools/sysmacros.h
		COMMAND sed -i -e "s/#include\ <features.h>//" kernel/linux-2.6.35/scripts/mod/elf.h
		COMMAND sed -i -e "s/#include\ <features.h>//" kernel/linux-2.6.35/arch/x86/boot/tools/sysmacros.h
		COMMAND sed -i -e "s,<sys/sysmacros.h>,\"sysmacros.h\"," kernel/linux-2.6.35/arch/x86/boot/tools/build.c
		COMMAND sed -i -e "s/<elf.h>/\"elf.h\"/" kernel/linux-2.6.35/scripts/mod/mk_elfconfig.c
		COMMAND sed -i -e "s/<elf.h>/\"elf.h\"/" kernel/linux-2.6.35/scripts/mod/modpost.h
	CONFIGURE_COMMAND	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS 			toolchain
)

ExternalProject_Add_Step(
	kernel
	hg_kernel
	DEPENDEES download
	DEPENDERS patch
	COMMAND hg clone https://code.google.com/p/googletv-mirrored-source.kernel <BINARY_DIR>/linux-2.6.35
)

ExternalProject_Add(
	osal
	URL 				${CMAKE_SOURCE_DIR}/osal.tar.gz
	SOURCE_DIR			osal
	BINARY_DIR			osal/osal
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS 			smd_common kernel
)

ExternalProject_Add(
	auto_eas
	URL 				${CMAKE_SOURCE_DIR}/auto_eas.tar.gz
	SOURCE_DIR 			auto_eas
	BINARY_DIR 			auto_eas/auto_eas
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal
)

ExternalProject_Add(
	htuple
	URL 				${CMAKE_SOURCE_DIR}/htuple.tar.gz
	SOURCE_DIR 			htuple
	BINARY_DIR 			htuple/htuple
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel
)

ExternalProject_Add(
	platform_config
	URL 				${CMAKE_SOURCE_DIR}/platform_config.tar.gz
	SOURCE_DIR 			platform_config
	BINARY_DIR 			platform_config/platform_config
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS htuple
)

ExternalProject_Add(
	pal
	URL					${CMAKE_SOURCE_DIR}/pal.tar.gz
	SOURCE_DIR 			pal
	BINARY_DIR 			pal/pal
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS auto_eas platform_config
)

ExternalProject_Add(
	sven
	URL 				${CMAKE_SOURCE_DIR}/sven.tar.gz
	SOURCE_DIR 			sven
	BINARY_DIR 			sven/sven
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS auto_eas platform_config pal
)

ExternalProject_Add(
	api
	URL 				${CMAKE_SOURCE_DIR}/api.tar.gz
	SOURCE_DIR 			api
	BINARY_DIR 			api/api
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	smd_tools
	URL 				${CMAKE_SOURCE_DIR}/smd_tools.tar.gz
	SOURCE_DIR 			smd_tools
	BINARY_DIR 			smd_tools/smd_tools
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	idl
	URL 				${CMAKE_SOURCE_DIR}/idl.tar.gz
	SOURCE_DIR 			idl
	BINARY_DIR 			idl/idl
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE} install
	DEPENDS osal sven
)

ExternalProject_Add(
	clock
	URL 				${CMAKE_SOURCE_DIR}/clock.tar.gz
	SOURCE_DIR 			clock
	BINARY_DIR 			clock/clock
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)


ExternalProject_Add(
	clock_control
	URL 				${CMAKE_SOURCE_DIR}/clock_control.tar.gz
	SOURCE_DIR 			clock_control
	BINARY_DIR 			clock_control/clock_control
	PATCH_COMMAND       sed -i -e "s/rc$<SEMICOLON>/rc=0$<SEMICOLON>/" clock_control/src/kernel/external_clock.c
		COMMAND         sed -i -e "s/-pd//" clock_control/src/lib/Makefile
		COMMAND         sed -i -e "s/-pd//" clock_control/src/kernel/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS idl
)

ExternalProject_Add(
	intel_ce_pm
	URL 				${CMAKE_SOURCE_DIR}/intel_ce_pm.tar.gz
	SOURCE_DIR 			intel_ce_pm
	BINARY_DIR 			intel_ce_pm/intel_ce_pm
	PATCH_COMMAND       sed -i -e "s/user//" intel_ce_pm/Makefile
		COMMAND         sed -i -e "s/-pd//" intel_ce_pm/kernel/Makefile
		COMMAND         sed -i -e "s/-pd//" intel_ce_pm/kernel/lib/Makefile
		COMMAND         ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/init_intel_ce_pm intel_ce_pm/init_intel_ce_pm
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS clock_control
)

ExternalProject_Add(
	system_utils
	URL 				${CMAKE_SOURCE_DIR}/system_utils.tar.gz
	SOURCE_DIR 			system_utils
	BINARY_DIR 			system_utils/system_utils
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	core
	URL 				${CMAKE_SOURCE_DIR}/core.tar.gz
	SOURCE_DIR 			core
	BINARY_DIR 			core/core
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven api smd_tools intel_ce_pm system_utils
)

ExternalProject_Add(
	nandflash
	URL 				${CMAKE_SOURCE_DIR}/nandflash.tar.gz
	SOURCE_DIR 			nandflash
	BINARY_DIR 			nandflash/nandflash
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel
)

ExternalProject_Add(
	generic_timer
	URL 				${CMAKE_SOURCE_DIR}/generic_timer.tar.gz
	SOURCE_DIR 			generic_timer
	BINARY_DIR 			generic_timer/generic_timer
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel
)

ExternalProject_Add(
	display
	URL 				${CMAKE_SOURCE_DIR}/display.tar.gz
	SOURCE_DIR 			display
	BINARY_DIR 			display/display
	PATCH_COMMAND       sed -i -e "s/-pd//" display/src/lib/libvbd/Makefile
		COMMAND         sed -i -e "s/-pd//" display/src/lib/libabd/Makefile
		COMMAND         sed -i -e "s/-pd//" display/src/kernel/Makefile
		COMMAND         sed -i -E "s/\\star\\s.*//" display/src/lib/libvbd/Makefile
		COMMAND         sed -i -E "s/\\star\\s.*//" display/src/lib/libabd/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel generic_timer
)


ExternalProject_Add(
	audio
	URL 				${CMAKE_SOURCE_DIR}/audio.tar.gz
	SOURCE_DIR 			audio
	BINARY_DIR 			audio/audio
	PATCH_COMMAND       sed -i -e "123,125d" audio/core/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	clock_recovery
	URL 				${CMAKE_SOURCE_DIR}/clock_recovery.tar.gz
	SOURCE_DIR 			clock_recovery
	BINARY_DIR 			clock_recovery/clock_recovery
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)
ExternalProject_Add(
	demux
	URL 				${CMAKE_SOURCE_DIR}/demux-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			demux
	BINARY_DIR 			demux/IntelCE/demux-13.11.10484.229643
	PATCH_COMMAND		sed -i -E "1251s/discon_detection_enabled\\s/discon_detection_enabled[stream_id]/" IntelCE/demux-13.11.10484.229643/ver3/driver/ismd_demux_tags_timing.c
		COMMAND			sed -i -e "4692s/$<SEMICOLON>/=odd_or_even$<SEMICOLON>/" IntelCE/demux-13.11.10484.229643/ver3/driver/ismd_demux_api.c
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core clock_recovery
)

ExternalProject_Add(
	viddec_fw
	URL 				${CMAKE_SOURCE_DIR}/viddec_fw-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			viddec_fw
	BINARY_DIR 			viddec_fw/IntelCE/viddec_fw-13.11.10504.232647
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	viddec
	URL 				${CMAKE_SOURCE_DIR}/viddec-DUAL-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			viddec
	BINARY_DIR 			viddec/IntelCE/viddec-13.11.10504.232647
	PATCH_COMMAND		sed -i -e "242d" -e "252d" -e "264,268d" IntelCE/viddec-13.11.10504.232647/src/ismd_viddec_pvt.h
		COMMAND 		sed -i -e "238d" -e "248d" -e "260,265d" IntelCE/viddec-13.11.10504.232647/version2/src/ismd_viddec_pvt.h
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS viddec_fw
)

ExternalProject_Add(
	common
	URL 				${CMAKE_SOURCE_DIR}/common.tar.gz
	SOURCE_DIR 			common
	BINARY_DIR 			common/common
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	ipclib
	URL 				${CMAKE_SOURCE_DIR}/ipclib-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			ipclib
	BINARY_DIR 			ipclib/IntelCE/ipclib-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	vidpproc_fw
	URL 				${CMAKE_SOURCE_DIR}/vidpproc_fw-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			vidpproc_fw
	BINARY_DIR 			vidpproc_fw/IntelCE/vidpproc_fw-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core ipclib
)

ExternalProject_Add(
	vidpproc
	URL 				${CMAKE_SOURCE_DIR}/vidpproc-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			vidpproc
	BINARY_DIR 			vidpproc/IntelCE/vidpproc-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS vidpproc_fw common
)

ExternalProject_Add(
	vidrend
	URL 				${CMAKE_SOURCE_DIR}/vidrend-None-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			vidrend
	BINARY_DIR 			vidrend/IntelCE/vidrend-13.11.10504.232647
	PATCH_COMMAND		sed -i -e "419,421d" -e "5774,5820d" -e "5823,5857d" IntelCE/vidrend-13.11.10504.232647/ismd_vidrend.c
		COMMAND			sed -i -e "291,293d" IntelCE/vidrend-13.11.10504.232647/vidrend.h
		COMMAND			sed -i -e "849,856d" -e "859d" IntelCE/vidrend-13.11.10504.232647/vidrend_th.c
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	tsout
	URL 				${CMAKE_SOURCE_DIR}/tsout-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			tsout
	BINARY_DIR 			tsout/IntelCE/tsout-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	avcap
	URL 				${CMAKE_SOURCE_DIR}/avcap-DUAL-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			avcap
	BINARY_DIR 			avcap/IntelCE/avcap-13.11.10504.232647
	PATCH_COMMAND		sed -i -e "62a #include <linux/slab.h>" IntelCE/avcap-13.11.10504.232647/src/kernel/core/avcap_core_queue.c
		COMMAND			sed -i -e "61a #include <linux/slab.h>" IntelCE/avcap-13.11.10504.232647/src/kernel/native_synthetic/synth_render.c
		COMMAND			sed -i -e "69d" IntelCE/avcap-13.11.10504.232647/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/lib/libvidcap_user/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/lib/libvidcap_kernel/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/kernel/core/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/kernel/native_synthetic/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/repeater/platform/custom/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/avcap-13.11.10504.232647/src/repeater/platform/internal/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)


ExternalProject_Add(
	bufmon
	URL 				${CMAKE_SOURCE_DIR}/bufmon-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			bufmon
	BINARY_DIR 			bufmon/IntelCE/bufmon-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS core
)

ExternalProject_Add(
	graphics
	URL 				${CMAKE_SOURCE_DIR}/graphics.tar.gz
	SOURCE_DIR 			graphics
	BINARY_DIR 			graphics/eurasia_km
	PATCH_COMMAND		${CMAKE_COMMAND} -E tar xvfz graphics/src/gfx_km_src.tgz
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		INTEL_D3_CHANGES=1 INSTALL_FSROOT=${BUILD_DEST} SUPPORT_SGX=1 SUPPORT_DRI_DRM=1 DISP_ROOT=${CMAKE_BINARY_DIR}/display/display GFX_ROOT=${CMAKE_BINARY_DIR}/display/display ${SMD_MAKE}
	DEPENDS kernel
)

#ExternalProject_Add(
#	sec
#	URL 				${CMAKE_SOURCE_DIR}/sec.tar.gz
#	SOURCE_DIR 			sec
#	BINARY_DIR 			sec/sec
#	CONFIGURE_COMMAND 	""
#	INSTALL_COMMAND		""
#	BUILD_COMMAND		${SMD_MAKE}
#	DEPENDS kernel
#)



