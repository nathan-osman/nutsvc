package monitor

import (
	"golang.org/x/sys/windows"
)

var (
	powrprof = windows.MustLoadDLL("PowrProf.dll")

	procSetSuspendState = powrprof.MustFindProc("SetSuspendState")

	constSE_SHUTDOWN_NAME = windows.StringToUTF16Ptr("SeShutdownPrivilege")
)

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
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

func shutdown() error {

	// Get the current process token
	var token windows.Token
	if err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&token,
	); err != nil {
		return err
	}

	// Get the LUID for the shutdown privilege
	tkp := windows.Tokenprivileges{}
	if err := windows.LookupPrivilegeValue(
		nil,
		constSE_SHUTDOWN_NAME,
		&tkp.Privileges[0].Luid,
	); err != nil {
		return err
	}
	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	// Get the shutdown privilege
	if err := windows.AdjustTokenPrivileges(
		token,
		false,
		&tkp,
		0,
		nil,
		nil,
	); err != nil {
		return err
	}

	// Initiate the shutdown
	if err := windows.InitiateSystemShutdownEx(
		windows.StringToUTF16Ptr(""),
		windows.StringToUTF16Ptr("UPS power lost"),
		0,
		true,
		false,
		windows.SHTDN_REASON_MAJOR_POWER|windows.SHTDN_REASON_MINOR_ENVIRONMENT,
	); err != nil {
		return err
	}

	return nil
}

func hibernate() error {
	return fnSetSuspendState(
		true,
		true,
	)
}
