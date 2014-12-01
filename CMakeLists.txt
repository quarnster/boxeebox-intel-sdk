cmake_minimum_required(VERSION 2.8)

include(ExternalProject)
find_program(GO go)
if(${GO-NOTFOUND})
	message(FATAL_ERROR "Go not found")
endif()
set(SYSROOT /Volumes/Skivavbild/boxeebox-xbmc/tools/boxeebox/sysroot)
set(BUILD_DEST ${CMAKE_BINARY_DIR}/target)
set(ENV{BUILD_DEST} ${BUILD_DEST})
set(ENV{KERNEL_VER} linux-2.6.28)
set(ENV{CROSS_COMPILE} ${BUILD_DEST}/bin/i686-cm-linux-)
set(ENV{ARCH} i386)
set(ENV{FSROOT} ${BUILD_DEST})
#message("export CROSS_COMPILE=$ENV{CROSS_COMPILE}")
#message("export KERNEL_VER=$ENV{KERNEL_VER}")
#message("export BUILD_DEST=$ENV{BUILD_DEST}")
#message("export ARCH=$ENV{ARCH}")

if(${VERBOSE})
	set(ENV{LOUD} 1)
endif()



add_custom_command(
	OUTPUT ${BUILD_DEST}/bin
	COMMAND ${CMAKE_COMMAND} -E make_directory ${BUILD_DEST}/bin
)

set(TARGETCMDS "FSROOT=$ENV{FSROOT}")
set(CMDS addr2line ar as c++ cpp g++ gcc gccbug gcov gprof ld nm objcopy objdump ranlib readelf size strings strip)
foreach(CMD ${CMDS})
	string(TOUPPER TARGET${CMD} t)
	string(REPLACE ++ "XX" t ${t})
	set(ENV{${t}}  $ENV{CROSS_COMPILE}${CMD})
	set(TARGETCMDS ${TARGETCMDS} "${t}=$ENV{${t}}")
	#message("export ${t}=$ENV{${t}}")
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

add_custom_target(
	audio_fw
	COMMAND ${GO} run ./extract_audio_fw.go
	WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
	DEPENDS ${CMAKE_SOURCE_DIR}/extract_audio_fw.go
)

ExternalProject_Add(
	smd_common
	URL 				${CMAKE_SOURCE_DIR}/SMD_Common-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			SMD_Common
	BINARY_DIR 			SMD_Common/IntelCE/SMD_Common-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS toolchain
)



ExternalProject_Add(
	kernel
	URL 				${CMAKE_SOURCE_DIR}/kernel-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR			kernel
	BINARY_DIR			kernel/IntelCE/kernel-13.11.10484.229643
	PATCH_COMMAND
		COMMAND ${CMAKE_COMMAND} -E chdir IntelCE/kernel-13.11.10484.229643/ ${CMAKE_COMMAND} -E tar xvfz linux-2.6.28-src.tar.gz
		COMMAND ${CMAKE_COMMAND} -E create_symlink /usr/include/malloc/malloc.h IntelCE/kernel-13.11.10484.229643/linux-2.6.28/scripts/genksyms/malloc.h
		COMMAND ${CMAKE_COMMAND} -E copy ${SYSROOT}/usr/include/elf.h IntelCE/kernel-13.11.10484.229643/linux-2.6.28/scripts/mod/elf.h
		COMMAND ${CMAKE_COMMAND} -E copy ${SYSROOT}/usr/include/sys/sysmacros.h IntelCE/kernel-13.11.10484.229643/linux-2.6.28/arch/x86/boot/tools/sysmacros.h
		COMMAND sed -i -e "s/#include\ <features.h>//" IntelCE/kernel-13.11.10484.229643/linux-2.6.28/scripts/mod/elf.h
		COMMAND sed -i -e "s/#include\ <features.h>//" IntelCE/kernel-13.11.10484.229643/linux-2.6.28/arch/x86/boot/tools/sysmacros.h
		COMMAND sed -i -e "s,<sys/sysmacros.h>,\"sysmacros.h\"," IntelCE/kernel-13.11.10484.229643/linux-2.6.28/arch/x86/boot/tools/build.c
		COMMAND sed -i -e "s/<elf.h>/\"elf.h\"/" IntelCE/kernel-13.11.10484.229643/linux-2.6.28/scripts/mod/mk_elfconfig.c
		COMMAND sed -i -e "s/<elf.h>/\"elf.h\"/" IntelCE/kernel-13.11.10484.229643/linux-2.6.28/scripts/mod/modpost.h
		COMMAND sed -i -E "s,make -C,make -j8 -C," IntelCE/kernel-13.11.10484.229643/build_all.sh
	CONFIGURE_COMMAND	""
	INSTALL_COMMAND		${CMAKE_COMMAND} -E create_symlink ${BUILD_DEST}/kernel/linux-2.6.28/include/linux ${BUILD_DEST}/kernel/linux-2.6.28/include/generated
		COMMAND			${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/pci-intelce_pm.h ${BUILD_DEST}/kernel/linux-2.6.28/include/linux
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS 			toolchain
)

# ExternalProject_Add_Step(
# 	kernel
# 	hg_kernel
# 	DEPENDEES download
# 	DEPENDERS patch
# 	COMMAND hg clone https://code.google.com/p/googletv-mirrored-source.kernel <BINARY_DIR>/linux-2.6.35
# )

ExternalProject_Add(
	osal
	URL 				${CMAKE_SOURCE_DIR}/osal-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR			osal
	BINARY_DIR			osal/IntelCE/osal-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS 			smd_common kernel
)

ExternalProject_Add(
	auto_eas
	URL 				${CMAKE_SOURCE_DIR}/auto_eas-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			auto_eas
	BINARY_DIR 			auto_eas/IntelCE/auto_eas-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal
)

ExternalProject_Add(
	htuple
	URL 				${CMAKE_SOURCE_DIR}/htuple-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			htuple
	BINARY_DIR 			htuple/IntelCE/htuple-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel smd_common
)

ExternalProject_Add(
	platform_config
	URL 				${CMAKE_SOURCE_DIR}/platform_config-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			platform_config
	BINARY_DIR 			platform_config/IntelCE/platform_config-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS htuple
)

ExternalProject_Add(
	pal
	URL					${CMAKE_SOURCE_DIR}/pal-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			pal
	BINARY_DIR 			pal/IntelCE/pal-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS auto_eas platform_config
)

ExternalProject_Add(
	sven
	URL 				${CMAKE_SOURCE_DIR}/sven-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			sven
	BINARY_DIR 			sven/IntelCE/sven-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS auto_eas platform_config pal
)

ExternalProject_Add(
	api
	URL 				${CMAKE_SOURCE_DIR}/api-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			api
	BINARY_DIR 			api/IntelCE/api-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	smd_tools
	URL 				${CMAKE_SOURCE_DIR}/smd_tools-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			smd_tools
	BINARY_DIR 			smd_tools/IntelCE/smd_tools-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	idl
	URL 				${CMAKE_SOURCE_DIR}/idl-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			idl
	BINARY_DIR 			idl/IntelCE/idl-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE} install
	DEPENDS osal sven
)

ExternalProject_Add(
	clock_control
	URL 				${CMAKE_SOURCE_DIR}/clock_control-DUAL-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			clock_control
	BINARY_DIR 			clock_control/IntelCE/clock_control-13.11.10484.229643
	PATCH_COMMAND
		COMMAND         sed -i -e "s/-pd//" IntelCE/clock_control-13.11.10484.229643/src/kernel/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS idl
)

ExternalProject_Add(
	clock
	URL 				${CMAKE_SOURCE_DIR}/clock-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			clock
	BINARY_DIR 			clock/IntelCE/clock-13.11.10484.229643
	PATCH_COMMAND		sed -i -e "55s/clock_control//" IntelCE/clock-13.11.10484.229643/src/src/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven clock_control
)



ExternalProject_Add(
	intel_ce_pm
	URL 				${CMAKE_SOURCE_DIR}/intel_ce_pm-DUAL-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			intel_ce_pm
	BINARY_DIR 			intel_ce_pm/IntelCE/intel_ce_pm-13.11.10484.229643
	PATCH_COMMAND		sed -i -e "63s/user//" -e "130s/doc//" IntelCE/intel_ce_pm-13.11.10484.229643/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/intel_ce_pm-13.11.10484.229643/kernel/Makefile
		COMMAND         sed -i -e "s/-pd//" IntelCE/intel_ce_pm-13.11.10484.229643/kernel/lib/Makefile
		COMMAND         ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/init_intel_ce_pm IntelCE/intel_ce_pm-13.11.10484.229643/init_intel_ce_pm
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS clock_control
)

ExternalProject_Add(
	system_utils
	URL 				${CMAKE_SOURCE_DIR}/system_utils-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			system_utils
	BINARY_DIR 			system_utils/IntelCE/system_utils-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven
)

ExternalProject_Add(
	core
	URL 				${CMAKE_SOURCE_DIR}/core-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			core
	BINARY_DIR 			core/IntelCE/core-13.11.10484.229643
	PATCH_COMMAND		sed -i -e "600s/from_tag_storage,/from_tag_storage=0,/" IntelCE/core-13.11.10484.229643/core/ismd_buffer_tags.c
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS osal sven api smd_tools intel_ce_pm system_utils
)

ExternalProject_Add(
	nandflash
	URL 				${CMAKE_SOURCE_DIR}/nandflash-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			nandflash
	BINARY_DIR 			nandflash/IntelCE/nandflash-13.11.10484.229643
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		PATH=$ENV{PATH}:${BUILD_DEST}/bin ${SMD_MAKE}
	DEPENDS kernel
)

# ExternalProject_Add(
# 	generic_timer
# 	URL 				${CMAKE_SOURCE_DIR}/generic_timer.tar.gz
# 	SOURCE_DIR 			generic_timer
# 	BINARY_DIR 			generic_timer/generic_timer
# 	CONFIGURE_COMMAND 	""
# 	INSTALL_COMMAND		""
# 	BUILD_COMMAND		${SMD_MAKE}
# 	DEPENDS kernel smd_tools
# )

ExternalProject_Add(
	display
	URL 				${CMAKE_SOURCE_DIR}/display-DUAL-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			display
	BINARY_DIR 			display/IntelCE/display-13.11.10504.232647
	PATCH_COMMAND		${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/gdl.h IntelCE/display-13.11.10504.232647/src/include
		COMMAND			${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/libgdl.h IntelCE/display-13.11.10504.232647/src/include
		COMMAND			sed -i -e "s/-pd//" IntelCE/display-13.11.10504.232647/src/lib/libvbd/Makefile
		COMMAND			sed -i -e "s/-pd//" IntelCE/display-13.11.10504.232647/src/lib/libabd/Makefile
		COMMAND			sed -i -e "s/-pd//" IntelCE/display-13.11.10504.232647/src/kernel/Makefile

	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS kernel idl system_utils intel_ce_pm
)



ExternalProject_Add(
	audio
	URL 				${CMAKE_SOURCE_DIR}/audio-None-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			audio
	BINARY_DIR 			audio/IntelCE/audio-13.11.10504.232647
	PATCH_COMMAND		sed -i -e "94a ifneq ($(strip $$(have_aud_inc)),)" -e "95a endif" IntelCE/audio-13.11.10504.232647/Makefile
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		CORE_A_NAME=dsp0 CORE_B_NAME=dsp1 ${SMD_MAKE}
	DEPENDS core display audio_fw
)

ExternalProject_Add(
	alsa_shim
	URL 				${CMAKE_SOURCE_DIR}/alsa_shim-None-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			alsa_shim
	BINARY_DIR 			alsa_shim/IntelCE/alsa_shim-13.11.10504.232647
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		KERNEL_BUILD_DIR=${BUILD_DEST}/kernel/linux-2.6.28 ${SMD_MAKE}
	DEPENDS audio
)

ExternalProject_Add(
	edl
	URL 				${CMAKE_SOURCE_DIR}/edl-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			edl
	BINARY_DIR 			edl/IntelCE/edl-13.11.10504.232647
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		PATH=$ENV{PATH}:${BUILD_DEST}/bin  ${SMD_MAKE}
	DEPENDS audio
)


ExternalProject_Add(
	clock_recovery
	URL 				${CMAKE_SOURCE_DIR}/clock_recovery-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			clock_recovery
	BINARY_DIR 			clock_recovery/IntelCE/clock_recovery-13.11.10484.229643
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
		COMMAND			sed -i -e "102s/clock_control//" IntelCE/demux-13.11.10484.229643/ver3/driver/Makefile
		COMMAND			sed -i -e "52s/clock_control//" IntelCE/demux-13.11.10484.229643/ver2/hal/tsi/Makefile
		COMMAND			sed -i -e "115s/clock_control//" IntelCE/demux-13.11.10484.229643/ver2/core/Makefile
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
	# PATCH_COMMAND		sed -i -e "242d" -e "252d" -e "264,268d" IntelCE/viddec-13.11.10504.232647/src/ismd_viddec_pvt.h
	# 	COMMAND 		sed -i -e "238d" -e "248d" -e "260,265d" IntelCE/viddec-13.11.10504.232647/version2/src/ismd_viddec_pvt.h
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		${SMD_MAKE}
	DEPENDS viddec_fw
)

ExternalProject_Add(
	common
	URL 				${CMAKE_SOURCE_DIR}/common-None-SRC-13.11.10484.229643.tgz
	SOURCE_DIR 			common
	BINARY_DIR 			common/IntelCE/common-13.11.10484.229643
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
	# PATCH_COMMAND		sed -i -e "419,421d" -e "5774,5820d" -e "5823,5857d" IntelCE/vidrend-13.11.10504.232647/ismd_vidrend.c
	# 	COMMAND			sed -i -e "291,293d" IntelCE/vidrend-13.11.10504.232647/vidrend.h
	# 	COMMAND			sed -i -e "849,856d" -e "859d" IntelCE/vidrend-13.11.10504.232647/vidrend_th.c
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
	DEPENDS core clock_recovery
)

ExternalProject_Add(
	graphics
	URL 				${CMAKE_SOURCE_DIR}/graphics-GPL-SRC-13.11.10504.232647.tgz
	SOURCE_DIR 			graphics
	BINARY_DIR 			graphics/eurasia_km
	PATCH_COMMAND		${CMAKE_COMMAND} -E tar xvfz IntelCE/graphics-13.11.10504.232647/src/gfx_km_src.tgz
		COMMAND			sed -i -E "s,SUPPORT_DRI_DRM = 1,SUPPORT_DRI_DRM = 0," eurasia_km/eurasiacon/build/linux/sgx_intel_ce/makefile.shared_conf
	CONFIGURE_COMMAND 	""
	INSTALL_COMMAND		""
	BUILD_COMMAND		INTEL_D3_CHANGES=1 INTEL_D3_PAD=1 INTEL_D3_NO_PCI_ENUM=1 INTEL_D3_DISABLE_TEXTURE_STREAM=1 INTEL_D3_P_CHANGES=1 INTEL_D3_CACHED_CBUF=1 INTEL_D3_FLUSH=1 INTEL_D3_PM=1 INSTALL_FSROOT=${BUILD_DEST} SUPPORT_SGX=1 DISP_ROOT=${CMAKE_BINARY_DIR}/display/IntelCE/display-13.11.10504.232647 GFX_ROOT=${CMAKE_BINARY_DIR}/display/IntelCE/display-13.11.10504.232647 ${SMD_MAKE}
	DEPENDS kernel osal idl intel_ce_pm
)

# ExternalProject_Add(
# 	sec
# 	URL 				${CMAKE_SOURCE_DIR}/sec-DUAL-SRC-13.11.10504.232647.tgz
# 	SOURCE_DIR 			sec
# 	BINARY_DIR 			sec/IntelCE/sec-13.11.10504.232647
# #	PATCH_COMMAND		${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/sec_types.h sec/kernel/sec_types.h
# 	CONFIGURE_COMMAND 	""
# 	INSTALL_COMMAND		""
# 	BUILD_COMMAND		${SMD_MAKE}
# 	DEPENDS kernel
# )



