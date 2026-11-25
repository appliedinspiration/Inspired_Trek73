# Cross-compiles Inspired Trek73 for every supported platform into .\dist.
# Run from anywhere; requires only the Go toolchain.
#   .\scripts\build-release.ps1 -Version 1.0.0
param([string]$Version = "dev")

$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$distPath = Join-Path $repoRoot "dist"
$versionPackage = "github.com/appliedinspiration/inspired_trek73/internal/version"
$date = Get-Date -Format "yyyy-MM-dd"

$commit = "unknown"
$gitCommand = Get-Command git -ErrorAction SilentlyContinue
if ($gitCommand) {
    $gitPath = $gitCommand.Source
    $gitOutput = & $gitPath -C $repoRoot rev-parse --short HEAD 2>$null
    if ($LASTEXITCODE -eq 0 -and -not [string]::IsNullOrWhiteSpace($gitOutput)) {
        $commit = $gitOutput.Trim()
    }
    else {
        Write-Warning "Unable to determine the Git commit; embedding 'unknown'."
    }
}
else {
    Write-Warning "Git is unavailable; embedding 'unknown' as the commit."
}

$ldflags = "-s -w -X $versionPackage.Version=$Version -X $versionPackage.BuildDate=$date -X $versionPackage.GitCommit=$commit"
$targets = @(
    @{ OS = "windows"; Arch = "amd64"; Name = "windows-amd64"; Extension = ".exe" },
    @{ OS = "linux";   Arch = "amd64"; Name = "linux-amd64";   Extension = "" },
    @{ OS = "linux";   Arch = "arm64"; Name = "linux-arm64";   Extension = "" },
    @{ OS = "linux";   Arch = "arm";   Arm = "7"; Name = "linux-armv7"; Extension = "" },
    @{ OS = "darwin";  Arch = "amd64"; Name = "darwin-amd64";  Extension = "" },
    @{ OS = "darwin";  Arch = "arm64"; Name = "darwin-arm64";  Extension = "" }
)
$environmentNames = @("CGO_ENABLED", "GOOS", "GOARCH", "GOARM")
$previousEnvironment = @{}
foreach ($name in $environmentNames) {
    $previousEnvironment[$name] = [Environment]::GetEnvironmentVariable($name, "Process")
}

try {
    Push-Location -LiteralPath $repoRoot
    try {
        New-Item -ItemType Directory -Path $distPath -Force | Out-Null

        foreach ($target in $targets) {
            $fileName = "inspired_trek73-$($target.Name)$($target.Extension)"
            $outputPath = Join-Path $distPath $fileName
            if (Test-Path -LiteralPath $outputPath) {
                Remove-Item -LiteralPath $outputPath -Force
            }
        }

        foreach ($target in $targets) {
            $fileName = "inspired_trek73-$($target.Name)$($target.Extension)"
            $outputPath = Join-Path $distPath $fileName
            $env:CGO_ENABLED = "0"
            $env:GOOS = $target.OS
            $env:GOARCH = $target.Arch
            if ($target.ContainsKey("Arm")) {
                $env:GOARM = $target.Arm
            }
            else {
                [Environment]::SetEnvironmentVariable("GOARM", $null, "Process")
            }

            & go build -trimpath -ldflags $ldflags -o $outputPath ./cmd/inspired_trek73
            if ($LASTEXITCODE -ne 0) {
                throw "Build failed for $fileName (exit code $LASTEXITCODE)."
            }
            if (-not (Test-Path -LiteralPath $outputPath -PathType Leaf)) {
                throw "Build reported success but did not produce $outputPath."
            }
            Write-Host "built $outputPath"
        }
    }
    finally {
        Pop-Location
    }
}
finally {
    foreach ($name in $environmentNames) {
        [Environment]::SetEnvironmentVariable($name, $previousEnvironment[$name], "Process")
    }
}
