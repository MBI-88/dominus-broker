# Paquete(s) a excluir del coverage
$excluded = 'tests/mocks'

Write-Host "`n🔍 Getting packages list...`n"

# Listar todos los paquetes y excluir los que coincidan con el patrón
$packages = go list ./... | Where-Object { $_ -notmatch $excluded }

Write-Host "✅ Packages included in coverage:`n"
$packages | ForEach-Object { Write-Host " - $_" }

# Unir los paquetes con comas para pasarlos a -coverpkg
$coverpkg = $packages -join ','

Write-Host "`n🚀 Running test...`n"

# Ejecutar go test y capturar salida línea por línea
$testOutput = go test -coverpkg="$coverpkg" -covermode=atomic -coverprofile="coverage.out" ./tests/... 2>&1

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

Write-Host "`nWhaching output`n"
go tool cover -func="coverage.out"