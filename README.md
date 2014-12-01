brew install doxygen

mkdir build
cd build
cmake ..
make -j8
# for me it errors out here after a while complaining about strlen, however it can't be that important as it doesn't happen when typing make again
make



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


