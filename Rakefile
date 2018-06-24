task :test do
    sh "go test -timeout 5s ./..."
end

task :build do
    sh "GOOS=linux GOARCH=amd64 go build -o gojikoServer.linux.amd64 goa/*.go"
    sh "GOOS=linux GOARCH=386 go build -o udpResponder.linux.386 udpResponder/*.go"
    sh "GOOS=windows GOARCH=amd64 go build -o gojikoServer.windows.amd64.exe goa/*.go"
    sh "GOOS=windows GOARCH=386 go build -o udpResponder.windows.386.exe udpResponder/*.go"
end
