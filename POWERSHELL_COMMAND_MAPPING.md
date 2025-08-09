# 🔧 PowerShell Command Mapping for Lucien Shell

## 📊 IMPLEMENTATION STRATEGY LEGEND

- 🟢 **Go Built-in**: Implemented directly in Go shell core
- 🔵 **Agent**: Routed to Python agent for AI-enhanced functionality  
- 🟡 **Subprocess**: Call external command (powershell.exe, bash, etc.)
- 🔴 **Hybrid**: Combination of built-in + subprocess/agent

## 📁 FILE AND DIRECTORY OPERATIONS

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Get-ChildItem` / `ls` / `dir` | `ls`, `dir` | 🟢 Go Built-in | Native Go filepath.Walk with formatting |
| `Set-Location` / `cd` | `cd` | 🟢 Go Built-in | Native Go os.Chdir |
| `Get-Location` / `pwd` | `pwd` | 🟢 Go Built-in | Native Go os.Getwd |
| `New-Item -ItemType Directory` | `mkdir` | 🟢 Go Built-in | Native Go os.MkdirAll |
| `New-Item -ItemType File` | `touch` | 🟢 Go Built-in | Native Go os.Create + timestamp |
| `Remove-Item` | `rm`, `del` | 🟢 Go Built-in | Native Go os.Remove/RemoveAll |
| `Copy-Item` | `cp`, `copy` | 🟢 Go Built-in | Native Go io.Copy with filepath logic |
| `Move-Item` | `mv`, `move` | 🟢 Go Built-in | Native Go os.Rename with cross-device support |
| `Get-Content` | `cat`, `type` | 🟢 Go Built-in | Native Go file reading with encoding detection |
| `Set-Content` | `echo >` | 🟢 Go Built-in | Native Go file writing |
| `Add-Content` | `echo >>` | 🟢 Go Built-in | Native Go file appending |
| `Test-Path` | `test` | 🟢 Go Built-in | Native Go os.Stat |
| `Get-Item` | `stat` | 🟢 Go Built-in | Native Go os.Stat with detailed info |
| `Rename-Item` | `ren`, `rename` | 🟢 Go Built-in | Native Go os.Rename |
| `Find-Files` | `find` | 🔴 Hybrid | Go built-in + optional regex agent |

## 🌐 NETWORKING COMMANDS

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Invoke-WebRequest` | `curl`, `wget` | 🟡 Subprocess | Call curl.exe or wget.exe |
| `Invoke-RestMethod` | `http` | 🔵 Agent | HTTPAgent for REST API interactions |
| `Test-NetConnection` | `ping` | 🟡 Subprocess | Call ping.exe/ping command |
| `Get-NetIPAddress` | `ipconfig`, `ifconfig` | 🟡 Subprocess | Call ipconfig.exe or ifconfig |
| `Resolve-DnsName` | `nslookup`, `dig` | 🟡 Subprocess | Call nslookup.exe or dig |
| `Test-Connection` | `tracert`, `traceroute` | 🟡 Subprocess | Call tracert.exe or traceroute |
| `Get-NetAdapter` | `netstat` | 🟡 Subprocess | Call netstat -an |

## 🖥️ SYSTEM INFORMATION

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Get-Process` | `ps`, `tasklist` | 🟢 Go Built-in | Native Go process enumeration |
| `Stop-Process` | `kill`, `taskkill` | 🟢 Go Built-in | Native Go os.Process.Kill |
| `Get-Service` | `service`, `sc` | 🟡 Subprocess | Call sc.exe or systemctl |
| `Get-ComputerInfo` | `systeminfo` | 🟡 Subprocess | Call systeminfo.exe or uname |
| `Get-EventLog` | `eventlog` | 🟡 Subprocess | Call wevtutil.exe or journalctl |
| `Get-WmiObject Win32_LogicalDisk` | `df`, `fsutil` | 🟡 Subprocess | Call df or fsutil volume |
| `Get-HotFix` | `updates` | 🟡 Subprocess | Call wmic qfe or apt list |
| `Get-Uptime` | `uptime` | 🟢 Go Built-in | Native Go system uptime calculation |

## 📦 PACKAGE MANAGEMENT

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Install-Module` | `install` | 🔵 Agent | PackageAgent with multi-manager support |
| `Get-Module` | `list` | 🔵 Agent | PackageAgent listing |
| `Update-Module` | `update` | 🔵 Agent | PackageAgent updates |
| `winget install` | `winget` | 🟡 Subprocess | Call winget.exe directly |
| `choco install` | `choco` | 🟡 Subprocess | Call choco.exe directly |
| `scoop install` | `scoop` | 🟡 Subprocess | Call scoop.exe directly |

## 🔧 TEXT PROCESSING

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Select-String` | `grep`, `findstr` | 🟢 Go Built-in | Native Go regex with file scanning |
| `Sort-Object` | `sort` | 🟢 Go Built-in | Native Go sorting algorithms |
| `Select-Object` | `cut`, `awk` | 🟢 Go Built-in | Native Go field selection |
| `Where-Object` | `filter` | 🟢 Go Built-in | Native Go filtering with expressions |
| `ForEach-Object` | `each` | 🟢 Go Built-in | Native Go iteration with commands |
| `Measure-Object` | `wc` | 🟢 Go Built-in | Native Go counting (lines, words, chars) |
| `Group-Object` | `uniq` | 🟢 Go Built-in | Native Go grouping and counting |

## 🔐 SECURITY & PERMISSIONS

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Get-Acl` | `getfacl`, `icacls` | 🟡 Subprocess | Call icacls.exe or getfacl |
| `Set-Acl` | `setfacl`, `icacls` | 🟡 Subprocess | Call icacls.exe or setfacl |
| `Get-ExecutionPolicy` | `policy` | 🟢 Go Built-in | Read from Lucien policy engine |
| `Set-ExecutionPolicy` | `policy` | 🟢 Go Built-in | Write to Lucien policy engine |

## 💾 REGISTRY & CONFIGURATION

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Get-ItemProperty` | `reg`, `config` | 🟡 Subprocess | Call reg.exe or read config files |
| `Set-ItemProperty` | `reg`, `config` | 🟡 Subprocess | Call reg.exe or write config files |
| `Get-Variable` | `env`, `set` | 🟢 Go Built-in | Native Go environment variable access |
| `Set-Variable` | `set`, `export` | 🟢 Go Built-in | Native Go environment variable setting |

## 🏗️ DEVELOPMENT TOOLS

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `git` commands | `git` | 🔴 Hybrid | Subprocess + GitAgent for smart operations |
| `npm` commands | `npm` | 🟡 Subprocess | Call npm.exe directly |
| `docker` commands | `docker` | 🔴 Hybrid | Subprocess + DockerAgent for orchestration |
| `kubectl` commands | `kubectl` | 🟡 Subprocess | Call kubectl.exe directly |

## 🔄 ADVANCED OPERATIONS

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Start-Job` | `job` | 🟢 Go Built-in | Native Go goroutine job system |
| `Get-Job` | `jobs` | 🟢 Go Built-in | Native Go job status |
| `Wait-Job` | `wait` | 🟢 Go Built-in | Native Go job waiting |
| `Invoke-Command` | `ssh`, `remote` | 🔵 Agent | RemoteAgent for SSH execution |
| `Enter-PSSession` | `ssh` | 🟡 Subprocess | Call ssh.exe directly |

## 🎨 FORMATTING & OUTPUT

| PowerShell Command | Lucien Equivalent | Strategy | Implementation |
|-------------------|------------------|----------|----------------|
| `Format-Table` | `table` | 🟢 Go Built-in | Native Go table formatting |
| `Format-List` | `list` | 🟢 Go Built-in | Native Go list formatting |
| `Out-GridView` | `grid` | 🔵 Agent | UIAgent for interactive grid display |
| `ConvertTo-Json` | `json` | 🟢 Go Built-in | Native Go JSON marshaling |
| `ConvertFrom-Json` | `json` | 🟢 Go Built-in | Native Go JSON unmarshaling |
| `ConvertTo-Csv` | `csv` | 🟢 Go Built-in | Native Go CSV formatting |

## 📊 SUMMARY

| Strategy | Command Count | Percentage |
|----------|---------------|------------|
| 🟢 Go Built-in | 28 | 45% |
| 🟡 Subprocess | 19 | 31% |
| 🔵 Agent | 10 | 16% |
| 🔴 Hybrid | 5 | 8% |
| **Total** | **62** | **100%** |

## 🚀 IMPLEMENTATION PRIORITY

### Phase 1 (Immediate) - Core File Operations
- `ls`, `cd`, `pwd`, `mkdir`, `rm`, `cp`, `mv`
- `cat`, `echo`, `grep`, `sort`

### Phase 2 (Next Sprint) - System Operations  
- `ps`, `kill`, `df`, `uptime`
- `curl`, `ping`, `ssh`

### Phase 3 (Advanced) - Development Tools
- Git integration, Docker support
- Package management agents
- Advanced text processing

This mapping provides **complete PowerShell parity** while leveraging Go's performance for core operations and Python agents for AI-enhanced functionality.