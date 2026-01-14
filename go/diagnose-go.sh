#!/bin/bash
# Go toolchain diagnostic script

LOG_FILE="c:/Users/Alexey/Code/biz-operating-system/stlph/.cursor/debug.log"

log_json() {
    local hypothesis_id=$1
    local location=$2
    local message=$3
    local data=$4
    local timestamp=$(date +%s%3N)
    
    echo "{\"sessionId\":\"debug-session\",\"runId\":\"run1\",\"hypothesisId\":\"$hypothesis_id\",\"location\":\"$location\",\"message\":\"$message\",\"data\":$data,\"timestamp\":$timestamp}" >> "$LOG_FILE"
}

# Hypothesis A: Check Go version
GO_VERSION=$(go version 2>&1)
log_json "A" "diagnose-go.sh:go-version" "Go version check" "{\"version\":\"$GO_VERSION\"}"

# Hypothesis B: Check GOROOT
GOROOT_VAL=${GOROOT:-$(go env GOROOT 2>/dev/null)}
if [ -n "$GOROOT_VAL" ]; then
    GOROOT_EXISTS="true"
    [ ! -d "$GOROOT_VAL" ] && GOROOT_EXISTS="false"
    STDLIB_PATH="$GOROOT_VAL/src"
    STDLIB_EXISTS="false"
    [ -d "$STDLIB_PATH" ] && STDLIB_EXISTS="true"
    log_json "B" "diagnose-go.sh:goroot" "GOROOT check" "{\"goroot\":\"$GOROOT_VAL\",\"gorootExists\":$GOROOT_EXISTS,\"stdlibPath\":\"$STDLIB_PATH\",\"stdlibExists\":$STDLIB_EXISTS}"
else
    log_json "B" "diagnose-go.sh:goroot" "GOROOT check" "{\"goroot\":\"\",\"gorootExists\":false,\"error\":\"GOROOT not set\"}"
fi

# Hypothesis C: Check GOPATH
GOPATH_VAL=${GOPATH:-$(go env GOPATH 2>/dev/null)}
if [ -z "$GOPATH_VAL" ]; then
    if [ "$OS" = "Windows_NT" ]; then
        GOPATH_VAL="$USERPROFILE/go"
    else
        GOPATH_VAL="$HOME/go"
    fi
fi
GOPATH_EXISTS="false"
[ -d "$GOPATH_VAL" ] && GOPATH_EXISTS="true"
log_json "C" "diagnose-go.sh:gopath" "GOPATH check" "{\"gopath\":\"$GOPATH_VAL\",\"gopathExists\":$GOPATH_EXISTS}"

# Hypothesis D: Check toolchain directory
if [ "$OS" = "Windows_NT" ]; then
    TOOLCHAIN_DIR="$LOCALAPPDATA/go/toolchain"
else
    TOOLCHAIN_DIR="$HOME/.cache/go-build"
fi
TOOLCHAIN_EXISTS="false"
[ -d "$TOOLCHAIN_DIR" ] && TOOLCHAIN_EXISTS="true"
log_json "D" "diagnose-go.sh:toolchain" "Toolchain directory check" "{\"toolchainDir\":\"$TOOLCHAIN_DIR\",\"toolchainDirExists\":$TOOLCHAIN_EXISTS}"

# Hypothesis E: Go env output
GO_ENV=$(go env GOROOT GOPATH GOTOOLCHAIN 2>&1)
log_json "E" "diagnose-go.sh:go-env" "Go environment variables" "{\"goEnvOutput\":\"$GO_ENV\"}"

# Hypothesis F: Check go.mod
if [ -f "go.mod" ]; then
    GO_MOD_FIRST_LINES=$(head -10 go.mod | tr '\n' ' ' | sed 's/"/\\"/g')
    GO_MOD_SIZE=$(wc -c < go.mod)
    log_json "F" "diagnose-go.sh:go-mod" "go.mod content check" "{\"goModExists\":true,\"goModSize\":$GO_MOD_SIZE,\"firstLines\":\"$GO_MOD_FIRST_LINES\"}"
else
    log_json "F" "diagnose-go.sh:go-mod" "go.mod content check" "{\"goModExists\":false}"
fi

# Hypothesis G: Try go list
GO_LIST_OUTPUT=$(cd go 2>/dev/null && go list -m 2>&1)
GO_LIST_EXIT=$?
log_json "G" "diagnose-go.sh:go-list" "go list command test" "{\"output\":\"$GO_LIST_OUTPUT\",\"exitCode\":$GO_LIST_EXIT}"

echo "Diagnostics complete. Check $LOG_FILE for details."
