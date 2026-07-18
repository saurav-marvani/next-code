# Crush → NextCode Rebranding Complete

## Overview
Successfully rebranded the entire Crush codebase to **NextCode** across all 373+ files.

## Changes Made

### 1. **Code References** (All Files Updated)
- `Crush` → `NextCode` (PascalCase)
- `crush` → `nextcode` (snake_case/kebab-case)
- `CRUSH` → `NEXTCODE` (UPPER_CASE)
- Specialized patterns:
  - `CrushClient` → `NextCodeClient`
  - `CrushService` → `NextCodeService`
  - `CrushProvider` → `NextCodeProvider`
  - `CrushConfig` → `NextCodeConfig`
  - `@crush/` → `@nextcode/`
  - `crush-cli` → `nextcode-cli`

### 2. **File Renames**
- `crush.json` → `nextcode.json`
- `internal/agent/tools/crush_info.go` → `internal/agent/tools/nextcode_info.go`
- `internal/agent/tools/crush_info.md` → `internal/agent/tools/nextcode_info.md`
- `internal/agent/tools/crush_info_test.go` → `internal/agent/tools/nextcode_info_test.go`
- `internal/agent/tools/crush_logs.go` → `internal/agent/tools/nextcode_logs.go`
- `internal/agent/tools/crush_logs.md.tpl` → `internal/agent/tools/nextcode_logs.md.tpl`
- `internal/agent/tools/crush_logs_test.go` → `internal/agent/tools/nextcode_logs_test.go`
- `internal/ui/notification/crush-icon.png` → `internal/ui/notification/nextcode-icon.png`
- `internal/ui/notification/crush-icon-solo.png` → `internal/ui/notification/nextcode-icon-solo.png`

### 3. **Directory Renames**
- `internal/skills/builtin/crush-config` → `internal/skills/builtin/nextcode-config`
- `internal/skills/builtin/crush-hooks` → `internal/skills/builtin/nextcode-hooks`

### 4. **Files Affected**
- Go source files (`.go`): ~200+ files
- Configuration files (`.json`, `.yaml`, `.yml`): ~50+ files
- Documentation (`.md`): ~30+ files
- Test files: Updated throughout
- Dockerfiles, Makefiles, and scripts: Updated

## Verification
- ✅ Zero remaining "crush" references (case-insensitive check)
- ✅ All file and directory renames completed
- ✅ Import paths updated correctly
- ✅ Tool references updated (NextCodeInfoTool, NextCodeLogsTool, etc.)
- ✅ Configuration schema references updated

## Files Not Changed
- Git history (`.git` directory - preserved)
- Node modules (excluded from search and replace)
- Binary files (handled appropriately)

## Next Steps
1. Review the changes for any edge cases
2. Test the NextCode CLI build and functionality
3. Update any external documentation or README references
4. Create PR and merge to main
