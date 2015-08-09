# Preparations

On a case sensitive file system (probably?):

	git clone --depth=1 https://github.com/quarnster/boxeebox-intel-sdk.git
	cd boxeebox-intel-sdk
	mkdir build
	cd build

# Prerequisites on darwin

	brew install doxygen gcc49 gnu-sed gettext gawk
	brew link --force gettext gmp4 libmpc08 mpfr2 isl011
	ulimit -n 1024
	ln -s /usr/local/bin/gsed sed
	export PATH=$PWD:$PATH
	export CC=/usr/local/bin/gcc-4.9
	export CXX=/usr/local/bin/g++-4.9
	export CPP=/usr/local/bin/cpp-4.9
	export LD=/usr/local/bin/gcc-4.9

# Building

	cmake ..
	make -j8


When the build fails, try just re-running "make" a couple of times and if the error persists, you'll probably have to tweak the CMakeLists.txt file a bit to suit your system (pull requests please) and re-run `make` a few times.


lib/modules/pvrsrvkm.ko
T SGXPostActivePowerEvent				      <
T SGXTestActivePowerEvent				      <
T graphics_pm_deinit					      <
T graphics_pm_init					      <
T graphics_pm_set_busy					      <
T graphics_pm_set_idle					      <
T graphics_pm_wait_not_suspended			      <
T icepm_device_register					      <
T icepm_device_unregister				      <
T icepm_set_drv_state					      <
T icepm_set_mode					      <
U intel_ce_pm_handler_v1				      <
U os_sema_destroy					      <
U os_sema_get						      <
U os_sema_init						      <
U os_sema_put						      <
U pal_flush_chipset_cache				      <
b graphics_pm_semaphore					      <
d graphics_pm_functions					      <
d graphics_pm_state					      <
t graphics_pm_resume					      <
t graphics_pm_suspend					      <


