@echo off
git pull && powershell -command "Stop-service -Force -name "DominatorRobot" -ErrorAction SilentlyContinue; go mod vendor; go build; Start-service -name "DominatorRobot""
:: Hail Hydra
