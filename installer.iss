#define MyAppName "Go PoC Windows Service"
#define MyAppVersion "0.0.1"
#define MyAppPublisher "Timo Reymann"
#define MyAppURL "https://github.com/timo-reymann/poc-go-windows-service"
#define MyAppExeName "windows-amd64.exe"

[Setup]
AppId={{1FEB2841-EDF0-43F7-8089-7BF74BC4454C}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
UninstallDisplayIcon={app}\{#MyAppExeName}
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
OutputBaseFilename=setup
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "dist\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"

[run]
Filename: {sys}\sc.exe; Parameters: "create go-poc-windows-service start= auto binPath= ""{app}\{#MyAppExeName}""" ; Flags: runhidden

[UninstallRun]
Filename: {sys}\sc.exe; Parameters: "stop go-poc-windows-service" ; Flags: runhidden
Filename: {sys}\sc.exe; Parameters: "delete go-poc-windows-service" ; Flags: runhidden
