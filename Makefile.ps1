#Requires -Version 5.1
<#
.SYNOPSIS
    Build and quality targets for dominus-project (Go module: dominus-project).

.DESCRIPTION
    PowerShell equivalent of a Makefile for Windows. Ensures commands run from the
    repository root. Go, golangci-lint, and govulncheck must be available on PATH
    (e.g. as configured by your Go version manager; this script does not alter PATH).

    All go test invocations use the race detector (-race) and -count=1 (no test-cache
    hits) so data races fail the run and are not masked by cached results.

    Coverage: -coverpkg from go list (internal minus boostrap + tests); see doc/coverage.md.

.PARAMETER Target
    Action to run (default: help).

.PARAMETER CoverageMinPct
    Minimum total coverage percentage for check-coverage (default matches param block below).

.PARAMETER AuditDocsDir
    Optional directory containing a release checklist (e.g. PENDIENTES.md). If empty,
    audit skips the documentation reminder.
#>

[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [ValidateSet(
        'build', 'test', 'test-cover', 'check-coverage',
        'fmt', 'lint', 'vuln', 'audit', 'deploy-check', 'help'
    )]
    [string] $Target = 'help',

    [ValidateRange(0, 100)]
    [int] $CoverageMinPct = 90,

    [string] $AuditDocsDir = ''
)

$ErrorActionPreference = 'Stop'

$repoRoot = $PSScriptRoot
Set-Location -LiteralPath $repoRoot

# Main package path and output binary name (see cmd/api/main.go)
$mainPackagePath = './cmd/api'
$binaryBaseName = 'dominus'

function Test-IsWindowsOS {
    return [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform(
        [System.Runtime.InteropServices.OSPlatform]::Windows)
}

function Get-BinaryFileName {
    if (Test-IsWindowsOS) { return "$binaryBaseName.exe" }
    return $binaryBaseName
}

function Write-StepMessage {
    param(
        [Parameter(Mandatory)]
        [string] $Message
    )
    Write-Host "[dominus-project] $Message" -ForegroundColor Cyan
}

function Assert-LastExitCode {
    param([string] $Step)
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed: $Step (exit code $LASTEXITCODE)" -ForegroundColor Red
        exit $LASTEXITCODE
    }
}

# Flags passed to every `go test` from this script (concurrency / race detection).
# -race: https://go.dev/doc/articles/race_detector — detects unsynchronized memory access; non-zero exit on race.
# -count=1: disables test result cache so a green cached run cannot skip a fresh -race execution.
function Get-GoTestConcurrencyFlags {
    return @('-race', '-count=1')
}

# Prints go test output: failures, panics, and race reports in red; ok lines in green.
function Write-GoTestOutputColored {
    param(
        [Parameter(Mandatory)]
        [AllowEmptyCollection()]
        [object[]] $Lines
    )
    foreach ($line in $Lines) {
        $text = [string]$line
        if ($text -match '(?i)(\bDATA\s+RACE\b|WARNING:\s*DATA\s+RACE|\brace detected at\b)') {
            Write-Host $text -ForegroundColor Red
        } elseif ($text -match '(?i)(^FAIL[\s\t]|^---\s*FAIL|\b---\s*FAIL:|^panic:|\bpanic:\s|testing\.Fatal|^\s*Error:\s|build failed|cannot find package|^\s*FAIL$)') {
            Write-Host $text -ForegroundColor Red
        } elseif ($text -match '^ok\s') {
            Write-Host $text -ForegroundColor Green
        } elseif ($text -match '^\?\s') {
            Write-Host $text -ForegroundColor DarkGray
        } else {
            Write-Host $text
        }
    }
}

function Invoke-GoTest {
    param(
        [Parameter(Mandatory)]
        [string[]] $TestArguments
    )
    $output = & go test @TestArguments 2>&1
    $exitCode = $LASTEXITCODE
    Write-GoTestOutputColored -Lines $output
    return $exitCode
}

function Get-AuditDocsPath {
    if ($AuditDocsDir -ne '') {
        return [IO.Path]::GetFullPath($AuditDocsDir)
    }
    return $null
}

function Invoke-BuildStep {
    $binDir = Join-Path $repoRoot 'bin'
    if (-not (Test-Path -LiteralPath $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
    $out = Join-Path $binDir (Get-BinaryFileName)
    Write-StepMessage "Building $mainPackagePath -> $out"
    & go build -o $out $mainPackagePath
    Assert-LastExitCode 'go build'
}

function Invoke-FmtStep {
    Write-StepMessage 'go fmt ./...'
    & go fmt ./...
    Assert-LastExitCode 'go fmt'
}

function Invoke-TestStep {
    Write-StepMessage 'go test ./... (race detector + -count=1)'
    $code = Invoke-GoTest -TestArguments (@('./...') + (Get-GoTestConcurrencyFlags))
    if ($code -ne 0) {
        Write-Host "Tests failed (exit code $code)." -ForegroundColor Red
        exit $code
    }
    Write-Host 'Successful!' -ForegroundColor Blue
}

function Get-CoverageProfilePath {
    return Join-Path $repoRoot 'coverage.out'
}

# Build -coverpkg as explicit import paths (go list). Excludes: cmd/, config/, mocks/
# (not listed), and internal/boostrap/ (bootstrap wiring, not unit-tested here).
function Get-CoveragePkgArg {
    $internalLines = & go list ./internal/...
    Assert-LastExitCode 'go list ./internal/...'
    $internalPkgs = foreach ($line in $internalLines) {
        $t = [string]$line.Trim()
        if ($t -eq '') { continue }
        if ($t -match '[/\\]boostrap$') { continue }
        $t
    }
    $testLines = & go list ./tests/...
    Assert-LastExitCode 'go list ./tests/...'
    $testPkgs = foreach ($line in $testLines) {
        $t = [string]$line.Trim()
        if ($t -eq '') { continue }
        $t
    }
    $all = @($internalPkgs) + @($testPkgs)
    if ($all.Count -eq 0) {
        throw 'Get-CoveragePkgArg: no packages from go list.'
    }
    return ($all -join ',')
}

function Invoke-TestCoverStep {
    param([switch] $SuppressSuccessBanner)

    $coverFile = Get-CoverageProfilePath
    $coverPkg = Get-CoveragePkgArg
    Write-StepMessage "Tests with coverage -> $coverFile (-coverpkg=$coverPkg)"
    $testArgs = @(
        "-coverprofile=$coverFile",
        "-coverpkg=$coverPkg"
    ) + (Get-GoTestConcurrencyFlags) + @('./...')
    $code = Invoke-GoTest -TestArguments $testArgs
    if ($code -ne 0) {
        Write-Host "Tests failed (exit code $code)." -ForegroundColor Red
        exit $code
    }
    if (-not (Test-Path -LiteralPath $coverFile)) {
        Write-Host 'No coverage profile was written.' -ForegroundColor Red
        exit 1
    }
    Write-Host '--- Coverage summary ---' -ForegroundColor DarkGray
    & go tool cover "-func=$coverFile"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed: go tool cover (exit code $LASTEXITCODE)" -ForegroundColor Red
        exit $LASTEXITCODE
    }
    if (-not $SuppressSuccessBanner) {
        Write-Host 'Successful!' -ForegroundColor Blue
    }
}

function Get-TotalCoveragePercent {
    $coverFile = Get-CoverageProfilePath
    if (-not (Test-Path -LiteralPath $coverFile)) {
        return $null
    }
    $coverReport = & go tool cover "-func=$coverFile" 2>&1
    if ($LASTEXITCODE -ne 0) {
        return $null
    }
    $lines = @(
        if ($null -eq $coverReport) {
            @()
        } elseif ($coverReport -is [System.Array]) {
            , $coverReport
        } else {
            , ($coverReport -split "`r?`n", [StringSplitOptions]::RemoveEmptyEntries)
        }
    )
    foreach ($line in $lines) {
        $text = [string]$line
        if ($text -notmatch '^\s*total:') {
            continue
        }
        if ($text -match '(\d+\.\d+|\d+)\s*%\s*$') {
            return [double]$Matches[1]
        }
    }
    return $null
}

function Invoke-CheckCoverageStep {
    Invoke-TestCoverStep -SuppressSuccessBanner
    $pct = Get-TotalCoveragePercent
    if ($null -eq $pct) {
        Write-Host 'Could not read total coverage.' -ForegroundColor Red
        exit 1
    }
    if ($pct -lt $CoverageMinPct) {
        Write-Host "Coverage $pct% is below minimum $CoverageMinPct%." -ForegroundColor Red
        exit 1
    }
    Write-Host "Coverage $pct% meets minimum $CoverageMinPct%." -ForegroundColor Green
    Write-Host 'Successful!' -ForegroundColor Blue
}

function Invoke-LintStep {
    Write-StepMessage 'golangci-lint run ./...'
    & golangci-lint run ./...
    Assert-LastExitCode 'golangci-lint'
}

function Invoke-VulnStep {
    Write-StepMessage 'govulncheck ./...'
    & govulncheck ./...
    Assert-LastExitCode 'govulncheck'
}

function Invoke-AuditStep {
    $docsRoot = Get-AuditDocsPath
    if ($null -ne $docsRoot) {
        $checklist = Join-Path $docsRoot 'PENDIENTES.md'
        if (Test-Path -LiteralPath $checklist) {
            Write-StepMessage "Release checklist: $checklist"
        } else {
            Write-Host "AuditDocsDir set but PENDIENTES.md not found: $checklist" -ForegroundColor Yellow
        }
    }
    Write-StepMessage 'go vet ./...'
    & go vet ./...
    Assert-LastExitCode 'go vet'
    Invoke-VulnStep
    Invoke-LintStep
}

function Invoke-DeployCheckStep {
    Write-StepMessage 'deploy-check: coverage + audit'
    Invoke-TestCoverStep -SuppressSuccessBanner
    Invoke-AuditStep
    Write-Host 'Successful! (deploy-check)' -ForegroundColor Blue
}

function Show-Help {
    $bin = Get-BinaryFileName
    Write-Host @"
dominus-project - local automation (PowerShell)

Usage:
  .\Makefile.ps1 [-Target <name>] [-CoverageMinPct <n>] [-AuditDocsDir <path>]

Targets:
  build           Build API binary to bin\$bin ($mainPackagePath)
  fmt             go fmt ./...
  test            go test ./... always with -race -count=1 (race detector, no test cache)
  test-cover      Same concurrency flags + coverage.out (coverpkg: internal + tests)
  check-coverage  test-cover then enforce -CoverageMinPct (see param default)
  lint            golangci-lint run ./...
  vuln            govulncheck ./...
  audit           go vet + vuln + lint; optional PENDIENTES.md if -AuditDocsDir is set
  deploy-check    test-cover + audit (pre-deploy gate)
  help            Show this text
"@ -ForegroundColor Gray
}

switch ($Target) {
    'build' { Invoke-BuildStep }
    'fmt' { Invoke-FmtStep }
    'test' { Invoke-TestStep }
    'test-cover' { Invoke-TestCoverStep }
    'check-coverage' { Invoke-CheckCoverageStep }
    'lint' { Invoke-LintStep }
    'vuln' { Invoke-VulnStep }
    'audit' { Invoke-AuditStep }
    'deploy-check' { Invoke-DeployCheckStep }
    'help' { Show-Help }
}
