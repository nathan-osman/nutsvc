package monitor

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	advapi32 = windows.MustLoadDLL("Advapi32.dll")
	powrprof = windows.MustLoadDLL("PowrProf.dll")

	procInitiateSystemShutdownExW = advapi32.MustFindProc("InitiateSystemShutdownExW")
	procSetSuspendState           = powrprof.MustFindProc("SetSuspendState")

	constSHTDN_REASON_MAJOR_POWER       = 0x00060000
	constSHTDN_REASON_MINOR_ENVIRONMENT = 0x0000000c
)

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func fnInitiateSystemShutdown(
	machineName, message string,
	timeout int,
	forceAppsClosed bool,
	rebootAfterShutdown bool,
	reason int,
) error {
	ret, _, err := procInitiateSystemShutdownExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(machineName))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(message))),
		uintptr(timeout),
		uintptr(boolToInt(forceAppsClosed)),
		uintptr(boolToInt(rebootAfterShutdown)),
		uintptr(reason),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func fnSetSuspendState(
	hibernate bool,
	wakeupEventsDisabled bool,
) error {
	ret, _, err := procSetSuspendState.Call(
		uintptr(boolToInt(hibernate)),
		uintptr(0),
		uintptr(boolToInt(wakeupEventsDisabled)),
	)
	if ret == 0 {
		return err
	}
	return nil
}
