; Windows x64 Installer for nutsvc
; Copyright 2024 - Nathan Osman

#define AppId        "{{a3500013-3e04-40f6-99fd-ea1ec65483cc}"
#define AppName      "nutsvc"
#define AppVersion   "0.1"
#define AppPublisher "Nathan Osman"
#define AppExe       "nutsvc.exe"

[Setup]
AppId={#AppId}
AppName={#AppName}
AppVersion={#AppVersion}
AppPublisher={#AppPublisher}
DefaultDirName={pf}\{#AppName}
DefaultGroupName={#AppName}
LicenseFile=LICENSE.txt
OutputDir=dist
OutputBaseFilename={#AppName}-{#AppVersion}-x86_64-setup
Compression=lzma
SolidCompression=yes
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64
CloseApplications=no

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: {#AppExe}; DestDir: "{app}"; BeforeInstall: PreInstall; AfterInstall: PostInstall

[Code]

// The service should only be installed if the executable didn't exist before
// installation. A global variable is required to keep track of this so that
// after installation, the appropriate action can be taken.

var
  ExeExisted: Boolean;

function ServiceCommand(Command: String): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{app}\{#AppExe}'), Command, '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
end;

procedure PreInstall();
begin
  ExeExisted := FileExists(ExpandConstant('{app}\{#AppExe}'));
  if ExeExisted then
  begin
    WizardForm.StatusLabel.Caption := 'Stopping service...';
    if not ServiceCommand('stop') then
      RaiseException('Unable to stop service.');
  end;
end;

procedure PostInstall();
begin
  if not ExeExisted then
  begin
    WizardForm.StatusLabel.Caption := 'Installing service...';
    if not ServiceCommand('install') then
      RaiseException('Unable to install service.');
  end;
  WizardForm.StatusLabel.Caption := 'Starting service...';
  if not ServiceCommand('start') then
      RaiseException('Unable to start service.');
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usUninstall then
  begin
    UninstallProgressForm.StatusLabel.Caption := 'Stopping & removing service...';
    if not ServiceCommand('stop') then
      RaiseException('Unable to stop service.');
    if not ServiceCommand('remove') then
      RaiseException('Unable to remove service.');
  end;
end;
