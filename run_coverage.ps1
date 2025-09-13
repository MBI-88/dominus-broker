# Lista de paquetes a excluir del coverage
$excluded = @(
    'tests/env',
    'docs',
    'mocks'
)

Write-Host "`n🔍 Getting packages list...`n"

# Listar todos los paquetes
$allPackages = go list ./...

# Filtrar los paquetes excluidos
$packages = $allPackages | Where-Object {
    $include = $true
    foreach ($pattern in $excluded) {
        if ($_ -match $pattern) {
            $include = $false
            break
        }
    }
    return $include
}

Write-Host "✅ Packages included in coverage:`n"
$packages | ForEach-Object { Write-Host " - $_" }

# Unir los paquetes con comas para pasarlos a -coverpkg
$coverpkg = $packages -join ','

Write-Host "`n🚀 Running test...`n"

# Ejecutar go test y capturar salida línea por línea
$testOutput = go test -timeout 120s -race -coverpkg="$coverpkg" -covermode=atomic -coverprofile="coverage.out" ./tests/... 2>&1

# Mostrar cada línea con formato y salto automático
foreach ($line in $testOutput) {
    if ($line -match 'FAIL') {
        Write-Host "$line `n" -ForegroundColor Red
    } elseif ($line -match 'PASS') {
        Write-Host "$line `n" -ForegroundColor Green
    } elseif ($line -match 'coverage:') {
        Write-Host "$line `n" -ForegroundColor Yellow
    } else {
        Write-Host "$line `n"
    }
}

Write-Host "`n📊 Watching output`n"
go tool cover -func="coverage.out"
