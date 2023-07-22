# Go build script for 32-bit and 64-bit architectures

# Set the name of the Go binary
$binaryName = "NorlandDemoUnlocker"

# Set the output directory for the compiled binaries
$outputDir = "./bin"

# Set the Go executable path
$goPath = "go"

# Set the target platforms (32-bit and 64-bit)
$targetPlatforms = @("386", "amd64")

# Loop through the target platforms and build the binaries
foreach ($platform in $targetPlatforms) {
    # Set the GOARCH environment variable
    $env:GOARCH = $platform

    $outputFile = Join-Path $outputDir "NorlandDemoUnlocker_$platform.exe"
    # Build the binary
    if ($platform -eq "386") {
        Write-Host "Building 32-bit binary..."
    } else {
        Write-Host "Building 64-bit binary..."
        $outputFile = Join-Path $outputDir "NorlandDemoUnlocker.exe"
    }
    & $goPath build -o $outputFile

    # Check if the build was successful
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Successfully built $platform-bit binary: $outputFile"
    } else {
        Write-Host "Build failed for $platform-bit binary."
        break
    }
}

# Clean up the GOARCH environment variable
Remove-Item Env:\GOARCH
