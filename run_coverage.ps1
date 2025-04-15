# Paquete(s) a excluir del coverage
$excluded = 'tests/mocks'

Write-Host "`n🔍 Obteniendo lista de paquetes...`n"

# Listar todos los paquetes y excluir los que coincidan con el patrón
$packages = go list ./... | Where-Object { $_ -notmatch $excluded }

Write-Host "✅ Paquetes incluidos en cobertura:`n"
$packages | ForEach-Object { Write-Host " - $_" }

# Unir los paquetes con comas para pasarlos a -coverpkg
$coverpkg = $packages -join ','

Write-Host "`n🚀 Ejecutando tests con cobertura...`n"

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

Write-Host "`n📄 Archivo de cobertura generado: go tool -func=coverage.out`n"
