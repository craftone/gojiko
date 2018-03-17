task :test do
    sh "go test -timeout 10s ./..."
end

task :build do
    sh "GOOS=linux GOARCH=amd64 go build -o gojikoServer.linux.amd64 goa/*.go"
    sh "GOOS=linux GOARCH=amd64 go build -o udpResponder.linux.amd64 udpResponder/*.go"
end
