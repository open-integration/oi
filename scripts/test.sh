go test -coverprofile cp.out ./pkg/modem
code=$?
go tool cover -html=cp.out -o coverage.html
echo "go test cmd exited with code $code"
exit $code