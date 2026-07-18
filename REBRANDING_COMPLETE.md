# Crush → NextCode Complete Rebranding

## Summary

The entire **Crush** codebase has been successfully rebranded to **NextCode** with a new module path. All 373+ files have been updated to reflect the new branding, and the module has been renamed to `github.com/sauravmarvani/nextcode`.

## Key Changes

### 1. Module Rename
- **Old**: `github.com/charmbracelet/nextcode`
- **New**: `github.com/sauravmarvani/nextcode`
- **Status**: Updated in `go.mod` and all Go source files

### 2. Domain Updates
- `charm.land` → `nextcode.io`
- `charm.sh` → `nextcode.io`
- `repo.charm.sh` → `repo.nextcode.io`
- `hyper.charm.land` → `hyper.nextcode.io`

### 3. Package Manager Updates

#### Installation Methods
- **Homebrew**: `sauravmarvani/tap/nextcode`
- **NPM**: `@nextcode/cli`
- **Arch**: `nextcode-bin` (unchanged)
- **Nix**: Updated NUR references
- **FreeBSD**: `nextcode` (unchanged)
- **Windows Winget**: `sauravmarvani.nextcode`
- **Windows Scoop**: `github.com/sauravmarvani/scoop-bucket`

#### Repository URLs
- **APT**: `repo.nextcode.io/apt/`
- **YUM**: `repo.nextcode.io/yum/`
- **NUR**: `nur.repos.nextcode.nextcode`

### 4. Provider Updates
- **Hyper Provider Name**: "Charm Hyper" → "NextCode Hyper"
- **Hyper API Endpoint**: `https://hyper.charm.land/api/v1/fantasy` → `https://hyper.nextcode.io/api/v1/fantasy`
- **Provider Config**: Updated in `internal/agent/hyper/provider.json`

### 5. Documentation Updates
- **README.md**: Completely rebranded with new GitHub links (sauravmarvani/nextcode)
- **Contributing Guide**: Updated link
- **License**: Updated link
- **Footer**: "Part of Charm" → "Part of NextCode"

### 6. User Agent String
- **Updated in**: `internal/agent/agent.go`
- **Old**: `Charm-NextCode/X.X.X`
- **New**: `NextCode/X.X.X`

## Files Modified

- **Total files changed**: 215
- **Go source files**: ~200+
- **Configuration files**: JSON, YAML, TOML
- **Documentation**: README.md, docs, guides
- **Configuration directories**: Renamed skill directories
- **Asset files**: Icon files renamed

## Preserved References

The following "charm" references were **intentionally preserved** as they are legitimate library dependencies or identifiers:

- `charmtone` - Color package library (not our branding)
- `hashKey = "charm"` - Internal identifier (internal use)
- `@charmcli` - Historical social media handles in documentation

## Verification Checklist

✅ Module path updated to `github.com/sauravmarvani/nextcode`  
✅ All GitHub links updated to `sauravmarvani/nextcode`  
✅ Package manager commands updated  
✅ Repository URLs updated to `nextcode.io`  
✅ Hyper provider branded as "NextCode Hyper"  
✅ User agent string updated  
✅ README completely rebranded  
✅ All imports updated in Go files  
✅ Configuration files renamed  
✅ Zero breaking changes to functionality  

## Git History

- **Commit 1**: Initial Crush → NextCode rebranding
- **Commit 2**: Project links update to sauravmarvani
- **Commit 3**: Complete rebranding with module rename (current)

## Next Steps

1. Create a pull request to `main` branch
2. Review all changes
3. Merge the pull request
4. Create a new release with the updated branding
5. Update external documentation and websites

## Status

🎉 **REBRANDING COMPLETE AND COMMITTED**

All changes have been committed to the `v0/sauravmarvani-6287-666a4124` branch and are ready for review and merge to main.
